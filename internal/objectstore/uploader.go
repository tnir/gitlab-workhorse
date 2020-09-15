package objectstore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
	"sync"
	"time"

	"gitlab.com/gitlab-org/labkit/log"
)

// uploader is an io.WriteCloser that can be used as write end of the uploading pipe.
type uploader struct {
	// etag is the object storage provided checksum
	etag string

	// md5 is an optional hasher for calculating md5 on the fly
	md5 hash.Hash

	w io.Writer

	// uploadError is the last error occourred during upload
	uploadError error
	// ctx is the internal context bound to the upload request
	ctx context.Context

	pr       *io.PipeReader
	pw       *io.PipeWriter
	strategy uploadStrategy
	metrics  bool

	// closeOnce is used to prevent multiple calls to pw.Close
	// which may result to Close overriding the error set by CloseWithError
	// Bug fixed in v1.14: https://github.com/golang/go/commit/f45eb9ff3c96dfd951c65d112d033ed7b5e02432
	closeOnce sync.Once
}

func newUploader(strategy uploadStrategy) uploader {
	pr, pw := io.Pipe()
	return uploader{w: pw, pr: pr, pw: pw, strategy: strategy, metrics: true}
}

func newMD5Uploader(strategy uploadStrategy, metrics bool) uploader {
	pr, pw := io.Pipe()
	hasher := md5.New()
	mw := io.MultiWriter(pw, hasher)
	return uploader{w: mw, pr: pr, pw: pw, md5: hasher, strategy: strategy, metrics: metrics}
}

// Close implements the standard io.Closer interface: it closes the http client request.
// This method will also wait for the connection to terminate and return any error occurred during the upload
func (u *uploader) Close() error {
	var closeError error
	u.closeOnce.Do(func() {
		closeError = u.pw.Close()
	})
	if closeError != nil {
		return closeError
	}

	<-u.ctx.Done()

	if err := u.ctx.Err(); err == context.DeadlineExceeded {
		return err
	}

	return u.uploadError
}

func (u *uploader) CloseWithError(err error) error {
	u.closeOnce.Do(func() {
		u.pw.CloseWithError(err)
	})

	return nil
}

func (u *uploader) Write(p []byte) (int, error) {
	return u.w.Write(p)
}

func (u *uploader) md5Sum() string {
	if u.md5 == nil {
		return ""
	}

	checksum := u.md5.Sum(nil)
	return hex.EncodeToString(checksum)
}

// ETag returns the checksum of the uploaded object returned by the ObjectStorage provider via ETag Header.
// This method will wait until upload context is done before returning.
func (u *uploader) ETag() string {
	<-u.ctx.Done()

	return u.etag
}

func (u *uploader) Execute(ctx context.Context, deadline time.Time) {
	if u.metrics {
		objectStorageUploadsOpen.Inc()
	}
	uploadCtx, cancelFn := context.WithDeadline(ctx, deadline)
	u.ctx = uploadCtx

	if u.metrics {
		go u.trackUploadTime()
	}

	uploadDone := make(chan struct{})
	go u.cleanup(ctx, uploadDone)
	go func() {
		defer cancelFn()
		defer close(uploadDone)

		if u.metrics {
			defer objectStorageUploadsOpen.Dec()
		}
		defer func() {
			// This will be returned as error to the next write operation on the pipe
			u.pr.CloseWithError(u.uploadError)
		}()

		err := u.strategy.Upload(uploadCtx, u.pr)
		if err != nil {
			u.uploadError = err
			if u.metrics {
				objectStorageUploadRequestsRequestFailed.Inc()
			}
			return
		}

		u.etag = u.strategy.ETag()

		if u.md5 != nil {
			err := compareMD5(u.md5Sum(), u.etag)
			if err != nil {
				log.ContextLogger(ctx).WithError(err).Error("error comparing MD5 checksum")

				u.uploadError = err
				if u.metrics {
					objectStorageUploadRequestsRequestFailed.Inc()
				}
			}
		}
	}()
}

func (u *uploader) trackUploadTime() {
	started := time.Now()
	<-u.ctx.Done()

	if u.metrics {
		objectStorageUploadTime.Observe(time.Since(started).Seconds())
	}
}

func (u *uploader) cleanup(ctx context.Context, uploadDone chan struct{}) {
	// wait for the upload to finish
	<-u.ctx.Done()

	<-uploadDone
	if u.uploadError != nil {
		if u.metrics {
			objectStorageUploadRequestsRequestFailed.Inc()
		}
		u.strategy.Abort()
		return
	}

	// We have now successfully uploaded the file to object storage. Another
	// goroutine will hand off the object to gitlab-rails.
	<-ctx.Done()

	// gitlab-rails is now done with the object so it's time to delete it.
	u.strategy.Delete()
}

func compareMD5(local, remote string) error {
	if !strings.EqualFold(local, remote) {
		return fmt.Errorf("ETag mismatch. expected %q got %q", local, remote)
	}

	return nil
}
