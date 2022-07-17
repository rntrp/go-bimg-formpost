package rest

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
)

// Valid GIFs have a minimal size of 14 bytes
// See https://github.com/mathiasbynens/small
const minValidFileSize = 14

// Theoretical max dimensions according to its respective specs:
//
// BMP (uint32 in go/x/image/bmp): 4,294,967,295
//
// GIF: 65,535
//
// JPEG: 65,535 (65,500 for the libjpeg-turbo based software)
//
// PNG: 4,294,967,295
//
// TIFF (uint32 in go/x/image/tiff): 4,294,967,295
//
// libjpeg-turbo max of 65,500 pixels appears to be a good limit for the other image
// formats, since all known image viewers are getting problems with image dimensions
// higher than this value, or even at values far smaller than this.
const MaxImageDim = 65_500

func coerceWidth(dim string) (int, error) {
	errFormat := "Supported width range is [%d;%d], got %d"
	return parseInt(1, MaxImageDim, 256, dim, errFormat)
}

func coerceHeight(dim string) (int, error) {
	errFormat := "Supported height range is [%d;%d], got %d"
	return parseInt(1, MaxImageDim, 256, dim, errFormat)
}

func parseInt(min, max, def int, num, errFormat string) (int, error) {
	if len(num) == 0 {
		return def, nil
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return n, err
	} else if min > n || n > max {
		return n, fmt.Errorf(errFormat, min, max, n)
	}
	return n, nil
}

func coerceFormat(format string) (bimg.ImageType, error) {
	switch strings.ToLower(format) {
	case "", "jpeg", "jpg":
		return bimg.JPEG, nil
	case "png":
		return bimg.PNG, nil
	case "gif":
		return bimg.GIF, nil
	case "tif", "tiff":
		return bimg.TIFF, nil
	case "webp":
		return bimg.WEBP, nil
	case "heif", "heic":
		return bimg.HEIF, nil
	case "avif":
		return bimg.AVIF, nil
	default:
		return bimg.UNKNOWN, fmt.Errorf("unknown image format '%s'", format)
	}
}

func coerceContentLength(contentLength string) (int64, error) {
	return strconv.ParseInt(contentLength, 10, 64)
}

const maxMemoryBufferSize = int64(math.MaxInt64) - 1

func coerceMemoryBufferSize(memoryBufferSize int64) int64 {
	if memoryBufferSize < 0 || memoryBufferSize > maxMemoryBufferSize {
		return maxMemoryBufferSize
	}
	return memoryBufferSize
}

func getOutputFileMeta(format bimg.ImageType) (string, string) {
	switch format {
	case bimg.JPEG:
		return ".jpg", "image/jpeg"
	case bimg.PNG:
		return ".png", "image/png"
	case bimg.GIF:
		return ".gif", "image/gif"
	case bimg.TIFF:
		return ".tif", "image/tiff"
	case bimg.WEBP:
		return ".webp", "image/webp"
	case bimg.HEIF:
		return ".heic", "image/heic"
	case bimg.AVIF:
		return ".avif", "image/avif"
	default:
		return "", "application/octet-stream"
	}
}
