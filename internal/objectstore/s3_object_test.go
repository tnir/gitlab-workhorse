package objectstore_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/objectstore"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/objectstore/test"
)

func TestS3ObjectUpload(t *testing.T) {
	creds, config, sess, ts := test.SetupS3(t)
	defer ts.Close()

	deadline := time.Now().Add(testTimeout)
	tmpDir, err := ioutil.TempDir("", "workhorse-test-")
	require.NoError(t, err)
	defer os.Remove(tmpDir)

	objectName := filepath.Join(tmpDir, "s3-test-data")
	ctx, cancel := context.WithCancel(context.Background())

	object, err := objectstore.NewS3Object(ctx, objectName, creds, config, deadline)
	require.NoError(t, err)

	// copy data
	n, err := io.Copy(object, strings.NewReader(test.ObjectContent))
	require.NoError(t, err)
	require.Equal(t, test.ObjectSize, n, "Uploaded file mismatch")

	// close HTTP stream
	err = object.Close()
	require.NoError(t, err)

	test.S3ObjectExists(t, sess, config, objectName, test.ObjectContent)

	cancel()
	deleted := false

	retry(3, time.Second, func() error {
		if test.S3ObjectDoesNotExist(t, sess, config, objectName) {
			deleted = true
			return nil
		} else {
			return fmt.Errorf("file is still present, retrying")
		}
	})

	require.True(t, deleted)
}

func TestConcurrentS3ObjectUpload(t *testing.T) {
	creds, uploadsConfig, uploadsSession, uploadServer := test.SetupS3WithBucket(t, "uploads")
	defer uploadServer.Close()

	// This will return a separate S3 endpoint
	_, artifactsConfig, artifactsSession, artifactsServer := test.SetupS3WithBucket(t, "artifacts")
	defer artifactsServer.Close()

	deadline := time.Now().Add(testTimeout)
	tmpDir, err := ioutil.TempDir("", "workhorse-test-")
	require.NoError(t, err)
	defer os.Remove(tmpDir)

	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func(index int) {
			var sess *session.Session
			var config config.S3Config

			if index%2 == 0 {
				sess = uploadsSession
				config = uploadsConfig
			} else {
				sess = artifactsSession
				config = artifactsConfig
			}

			name := fmt.Sprintf("s3-test-data-%d", index)
			objectName := filepath.Join(tmpDir, name)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			object, err := objectstore.NewS3Object(ctx, objectName, creds, config, deadline)
			require.NoError(t, err)

			// copy data
			n, err := io.Copy(object, strings.NewReader(test.ObjectContent))
			require.NoError(t, err)
			require.Equal(t, test.ObjectSize, n, "Uploaded file mismatch")

			// close HTTP stream
			require.NoError(t, object.Close())

			test.S3ObjectExists(t, sess, config, objectName, test.ObjectContent)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestS3ObjectUploadCancel(t *testing.T) {
	creds, config, _, ts := test.SetupS3(t)
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())

	deadline := time.Now().Add(testTimeout)
	tmpDir, err := ioutil.TempDir("", "workhorse-test-")
	require.NoError(t, err)
	defer os.Remove(tmpDir)

	objectName := filepath.Join(tmpDir, "s3-test-data")

	object, err := objectstore.NewS3Object(ctx, objectName, creds, config, deadline)

	require.NoError(t, err)

	// Cancel the transfer before the data has been copied to ensure
	// we handle this gracefully.
	cancel()

	_, err = io.Copy(object, strings.NewReader(test.ObjectContent))
	require.Error(t, err)
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}