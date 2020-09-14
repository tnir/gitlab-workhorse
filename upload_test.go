package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/secret"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/testhelper"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/upload"
)

type uploadArtifactsFunction func(url, contentType string, body io.Reader) (*http.Response, string, error)

func uploadArtifactsV1(url, contentType string, body io.Reader) (*http.Response, string, error) {
	resource := `/ci/api/v1/builds/123/artifacts`
	resp, err := http.Post(url+resource, contentType, body)
	return resp, resource, err
}

func uploadArtifactsV4(url, contentType string, body io.Reader) (*http.Response, string, error) {
	resource := `/api/v4/jobs/123/artifacts`
	resp, err := http.Post(url+resource, contentType, body)
	return resp, resource, err
}

func testArtifactsUpload(t *testing.T, uploadArtifacts uploadArtifactsFunction) {
	reqBody, contentType, err := multipartBodyWithFile()
	require.NoError(t, err)

	ts := signedUploadTestServer(t, nil)
	defer ts.Close()

	ws := startWorkhorseServer(ts.URL)
	defer ws.Close()

	resp, resource, err := uploadArtifacts(ws.URL, contentType, reqBody)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "GET %q: expected 200, got %d", resource, resp.StatusCode)
}

func TestArtifactsUpload(t *testing.T) {
	testArtifactsUpload(t, uploadArtifactsV1)
	testArtifactsUpload(t, uploadArtifactsV4)
}

func expectSignedRequest(t *testing.T, r *http.Request) {
	t.Helper()

	_, err := jwt.Parse(r.Header.Get(secret.RequestHeader), testhelper.ParseJWT)
	require.NoError(t, err)
}

func uploadTestServer(t *testing.T, extraTests func(r *http.Request)) *httptest.Server {
	return testhelper.TestServerWithHandler(regexp.MustCompile(`.`), func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/authorize") {
			expectSignedRequest(t, r)

			w.Header().Set("Content-Type", api.ResponseContentType)
			_, err := fmt.Fprintf(w, `{"TempPath":"%s"}`, scratchDir)
			require.NoError(t, err)
			return
		}

		require.NoError(t, r.ParseMultipartForm(100000))

		const nValues = 11 // file name, path, remote_url, remote_id, size, md5, sha1, sha256, sha512, gitlab-workhorse-upload, etag for just the upload (no metadata because we are not POSTing a valid zip file)
		require.Len(t, r.MultipartForm.Value, nValues)

		require.Empty(t, r.MultipartForm.File, "multipart form files")

		if extraTests != nil {
			extraTests(r)
		}
		w.WriteHeader(200)
	})
}

func signedUploadTestServer(t *testing.T, extraTests func(r *http.Request)) *httptest.Server {
	t.Helper()

	return uploadTestServer(t, func(r *http.Request) {
		expectSignedRequest(t, r)

		if extraTests != nil {
			extraTests(r)
		}
	})
}

func TestAcceleratedUpload(t *testing.T) {
	tests := []struct {
		method             string
		resource           string
		signedFinalization bool
	}{
		{"POST", `/example`, false},
		{"POST", `/uploads/personal_snippet`, true},
		{"POST", `/uploads/user`, true},
		{"POST", `/api/v4/projects/1/wikis/attachments`, false},
		{"POST", `/api/graphql`, false},
		{"PUT", "/api/v4/projects/9001/packages/nuget/v1/files", true},
		{"POST", `/api/v4/groups/import`, true},
		{"POST", `/api/v4/projects/import`, true},
		{"POST", `/import/gitlab_project`, true},
		{"POST", `/import/gitlab_group`, true},
		{"POST", `/api/v4/projects/9001/packages/pypi`, true},
	}

	for _, tt := range tests {
		t.Run(tt.resource, func(t *testing.T) {
			ts := uploadTestServer(t,
				func(r *http.Request) {
					if tt.signedFinalization {
						expectSignedRequest(t, r)
					}

					token, err := jwt.ParseWithClaims(r.Header.Get(upload.RewrittenFieldsHeader), &upload.MultipartClaims{}, testhelper.ParseJWT)
					require.NoError(t, err)

					rewrittenFields := token.Claims.(*upload.MultipartClaims).RewrittenFields
					if len(rewrittenFields) != 1 || len(rewrittenFields["file"]) == 0 {
						t.Fatalf("Unexpected rewritten_fields value: %v", rewrittenFields)
					}

					token, jwtErr := jwt.ParseWithClaims(r.PostFormValue("file.gitlab-workhorse-upload"), &testhelper.UploadClaims{}, testhelper.ParseJWT)
					require.NoError(t, jwtErr)

					uploadFields := token.Claims.(*testhelper.UploadClaims).Upload
					require.Contains(t, uploadFields, "name")
					require.Contains(t, uploadFields, "path")
					require.Contains(t, uploadFields, "remote_url")
					require.Contains(t, uploadFields, "remote_id")
					require.Contains(t, uploadFields, "size")
					require.Contains(t, uploadFields, "md5")
					require.Contains(t, uploadFields, "sha1")
					require.Contains(t, uploadFields, "sha256")
					require.Contains(t, uploadFields, "sha512")
				})

			defer ts.Close()
			ws := startWorkhorseServer(ts.URL)
			defer ws.Close()

			reqBody, contentType, err := multipartBodyWithFile()
			require.NoError(t, err)

			req, err := http.NewRequest(tt.method, ws.URL+tt.resource, reqBody)
			require.NoError(t, err)

			req.Header.Set("Content-Type", contentType)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			resp.Body.Close()
		})
	}
}

func multipartBodyWithFile() (io.Reader, string, error) {
	result := &bytes.Buffer{}
	writer := multipart.NewWriter(result)
	file, err := writer.CreateFormFile("file", "my.file")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprint(file, "SHOULD BE ON DISK, NOT IN MULTIPART")
	return result, writer.FormDataContentType(), writer.Close()
}

func TestBlockingRewrittenFieldsHeader(t *testing.T) {
	canary := "untrusted header passed by user"
	testCases := []struct {
		desc        string
		contentType string
		body        io.Reader
		present     bool
	}{
		{"multipart with file", "", nil, true}, // placeholder
		{"no multipart", "text/plain", nil, false},
	}

	var err error
	testCases[0].body, testCases[0].contentType, err = multipartBodyWithFile()
	require.NoError(t, err)

	for _, tc := range testCases {
		ts := testhelper.TestServerWithHandler(regexp.MustCompile(`.`), func(w http.ResponseWriter, r *http.Request) {
			key := upload.RewrittenFieldsHeader
			if tc.present {
				require.Contains(t, r.Header, key)
			} else {
				require.NotContains(t, r.Header, key)
			}

			require.NotEqual(t, canary, r.Header.Get(key), "Found canary %q in header %q", canary, key)
		})
		defer ts.Close()
		ws := startWorkhorseServer(ts.URL)
		defer ws.Close()

		req, err := http.NewRequest("POST", ws.URL+"/something", tc.body)
		require.NoError(t, err)

		req.Header.Set("Content-Type", tc.contentType)
		req.Header.Set(upload.RewrittenFieldsHeader, canary)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode, "status code")
	}
}

func TestLfsUpload(t *testing.T) {
	reqBody := "test data"
	rspBody := "test success"
	oid := "916f0027a575074ce72a331777c3478d6513f786a591bd892da1a577bf2335f9"
	resource := fmt.Sprintf("/%s/gitlab-lfs/objects/%s/%d", testRepo, oid, len(reqBody))

	lfsApiResponse := fmt.Sprintf(
		`{"TempPath":%q, "LfsOid":%q, "LfsSize": %d}`,
		scratchDir, oid, len(reqBody),
	)

	ts := testhelper.TestServerWithHandler(regexp.MustCompile(`.`), func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, "PUT")
		switch r.RequestURI {
		case resource + "/authorize":
			expectSignedRequest(t, r)

			// Instruct workhorse to accept the upload
			w.Header().Set("Content-Type", api.ResponseContentType)
			_, err := fmt.Fprint(w, lfsApiResponse)
			require.NoError(t, err)

		case resource:
			expectSignedRequest(t, r)

			// Expect the request to point to a file on disk containing the data
			require.NoError(t, r.ParseForm())
			require.Equal(t, oid, r.Form.Get("file.sha256"), "Invalid SHA256 populated")
			require.Equal(t, strconv.Itoa(len(reqBody)), r.Form.Get("file.size"), "Invalid size populated")

			tempfile, err := ioutil.ReadFile(r.Form.Get("file.path"))
			require.NoError(t, err)
			require.Equal(t, reqBody, string(tempfile), "Temporary file has the wrong body")

			fmt.Fprint(w, rspBody)
		default:
			t.Fatalf("Unexpected request to upstream! %v %q", r.Method, r.RequestURI)
		}
	})
	defer ts.Close()

	ws := startWorkhorseServer(ts.URL)
	defer ws.Close()

	req, err := http.NewRequest("PUT", ws.URL+resource, strings.NewReader(reqBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = int64(len(reqBody))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	rspData, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	// Expect the (eventual) response to be proxied through, untouched
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, rspBody, string(rspData))
}

func packageUploadTestServer(t *testing.T, resource string, reqBody string, rspBody string) *httptest.Server {
	return testhelper.TestServerWithHandler(regexp.MustCompile(`.`), func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, "PUT")
		apiResponse := fmt.Sprintf(
			`{"TempPath":%q, "Size": %d}`, scratchDir, len(reqBody),
		)
		switch r.RequestURI {
		case resource + "/authorize":
			expectSignedRequest(t, r)

			// Instruct workhorse to accept the upload
			w.Header().Set("Content-Type", api.ResponseContentType)
			_, err := fmt.Fprint(w, apiResponse)
			require.NoError(t, err)

		case resource:
			expectSignedRequest(t, r)

			// Expect the request to point to a file on disk containing the data
			require.NoError(t, r.ParseForm())

			len := strconv.Itoa(len(reqBody))
			require.Equal(t, len, r.Form.Get("file.size"), "Invalid size populated")

			tmpFilePath := r.Form.Get("file.path")
			fileData, err := ioutil.ReadFile(tmpFilePath)
			defer os.Remove(tmpFilePath)

			require.NoError(t, err)
			require.Equal(t, reqBody, string(fileData), "Temporary file has the wrong body")

			fmt.Fprint(w, rspBody)
		default:
			t.Fatalf("Unexpected request to upstream! %v %q", r.Method, r.RequestURI)
		}
	})
}

func testPackageFileUpload(t *testing.T, resource string) {
	reqBody := "test data"
	rspBody := "test success"

	ts := packageUploadTestServer(t, resource, reqBody, rspBody)
	defer ts.Close()

	ws := startWorkhorseServer(ts.URL)
	defer ws.Close()

	req, err := http.NewRequest("PUT", ws.URL+resource, strings.NewReader(reqBody))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respData, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, rspBody, string(respData), "Temporary file has the wrong body")
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}

func TestPackageFilesUpload(t *testing.T) {
	routes := []string{
		"/api/v4/packages/conan/v1/files",
		"/api/v4/projects/2412/packages/conan/v1/files",
		"/api/v4/projects/2412/packages/maven/v1/files",
	}

	for _, r := range routes {
		testPackageFileUpload(t, r)
	}
}
