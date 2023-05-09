package rest

import (
	"encoding/hex"
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
	errFormat := "supported width range is [%d;%d], got %d"
	return parseInt(1, MaxImageDim, 256, dim, errFormat)
}

func coerceHeight(dim string) (int, error) {
	errFormat := "supported height range is [%d;%d], got %d"
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

func coerceQuality(quality string, format bimg.ImageType) (int, int, error) {
	switch format {
	case bimg.JPEG:
		q, err := parseInt(0, 100, 99, quality,
			"supported JPEG quality range is range is [%d;%d], got %d")
		return 0, q, err
	case bimg.PNG:
		c, err := parseInt(0, 9, 6, quality,
			"supported PNG compression level range is [%d;%d], got %d")
		return c, 0, err
	default:
		return 0, 0, nil
	}
}

func coerceResample(resample string) (bimg.Interpolator, error) {
	switch strings.ToLower(resample) {
	case "", "bicubic":
		return bimg.Bicubic, nil
	case "bilinear":
		return bimg.Bilinear, nil
	case "nohalo":
		return bimg.Nohalo, nil
	case "nearest":
		return bimg.Nearest, nil
	default:
		return bimg.Bicubic, fmt.Errorf("unknown resampling mode: %s", resample)
	}
}

func coerceBackground(background string) (bimg.Color, error) {
	if len(background) == 0 {
		return bimg.ColorBlack, nil
	}
	color, err := hex.DecodeString(strings.TrimPrefix(background, "#"))
	if err != nil {
		return bimg.ColorBlack, err
	} else if len(color) != 3 {
		return bimg.ColorBlack, fmt.Errorf("invalid RGB hex string: %s", background)
	}
	return bimg.Color{R: color[0], G: color[1], B: color[2]}, nil
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

func coerceResize(resize string, options *bimg.Options) error {
	switch strings.ToLower(resize) {
	case "", "fit":
		options.Crop = false
		options.Embed = true
		options.Enlarge = false
		options.Force = false
	case "fit-upscale":
		options.Crop = false
		options.Embed = true
		options.Enlarge = true
		options.Force = true
	case "fit-upscale-black":
		options.Crop = false
		options.Embed = true
		options.Enlarge = true
		options.Extend = bimg.ExtendBlack
		options.Force = false
	case "fit-upscale-white":
		options.Crop = false
		options.Embed = true
		options.Enlarge = true
		options.Extend = bimg.ExtendWhite
		options.Force = false
	case "fill":
		coerceResizeFill(options, bimg.GravityCentre)
	case "fill-north":
		coerceResizeFill(options, bimg.GravityNorth)
	case "fill-east":
		coerceResizeFill(options, bimg.GravityEast)
	case "fill-south":
		coerceResizeFill(options, bimg.GravitySouth)
	case "fill-west":
		coerceResizeFill(options, bimg.GravityWest)
	case "fill-smart":
		coerceResizeFill(options, bimg.GravitySmart)
	case "stretch":
		options.Crop = false
		options.Force = true
		options.Embed = true
		options.Enlarge = true
	default:
		return fmt.Errorf("unknown resize mode: %s", resize)
	}
	return nil
}

func coerceResizeFill(options *bimg.Options, gravity bimg.Gravity) {
	options.Crop = true
	options.Embed = true
	options.Enlarge = true
	options.Force = false
	options.Gravity = gravity
}

func preserveAspectRatio(resize string, src bimg.ImageSize, options *bimg.Options) {
	switch strings.ToLower(resize) {
	case "", "fit", "fit-upscale":
		wf := float64(options.Width) / float64(src.Width)
		hf := float64(options.Height) / float64(src.Height)
		if wf > hf {
			options.Width = int(math.Round(float64(src.Width) * hf))
		} else if wf < hf {
			options.Height = int(math.Round(float64(src.Height) * wf))
		}
	}
}
