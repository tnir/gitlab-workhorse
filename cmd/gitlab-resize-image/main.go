package main

import (
	"fmt"
	"image"
	"mime"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

func main() {
	widthParam := os.Getenv("GL_RESIZE_IMAGE_WIDTH")
	requestedWidth, err := strconv.Atoi(widthParam)
	if err != nil {
		fail("failed parsing GL_RESIZE_IMAGE_WIDTH; not a valid integer:", widthParam)
	}
	contentType := os.Getenv("GL_RESIZE_IMAGE_CONTENT_TYPE")
	if contentType == "" {
		fail("GL_RESIZE_IMAGE_CONTENT_TYPE not set")
	}

	src, extension, err := image.Decode(os.Stdin)
	if err != nil {
		fail("failed to open image:", err)
	}
	if detectedType := mime.TypeByExtension("." + extension); detectedType != contentType {
		fail("MIME types do not match; requested:", contentType, "actual:", detectedType)
	}
	format, err := imaging.FormatFromExtension(extension)
	if err != nil {
		fail("failed to find extension:", err)
	}

	image := imaging.Resize(src, requestedWidth, 0, imaging.Lanczos)
	imaging.Encode(os.Stdout, image, format)
}

func fail(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
