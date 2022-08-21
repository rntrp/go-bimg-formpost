package rest

import (
	"log"
	"net/http"

	"github.com/h2non/bimg"
)

func Convert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, errW := coerceWidth(query.Get("width"))
	height, errH := coerceHeight(query.Get("height"))
	format, errF := coerceFormat(query.Get("format"))
	compression, quality, errCQ := coerceQuality(query.Get("quality"), format)
	resample, errRS := coerceResample(query.Get("resample"))
	background, errBKG := coerceBackground(query.Get("background"))
	for _, err := range [...]error{errW, errH, errF, errCQ, errRS, errBKG} {
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
	}
	options := bimg.Options{
		Width:         width,
		Height:        height,
		Type:          format,
		Compression:   compression,
		Quality:       quality,
		Interpolator:  resample,
		Background:    background,
		StripMetadata: true,
		NoProfile:     true,
	}
	if err := coerceResize(query.Get("resize"), &options); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	src := handleFileUpload(w, r)
	if src == nil {
		return
	}
	img := bimg.NewImage(src)
	size, sizeOk := handleImageSize(w, img)
	if !sizeOk {
		return
	}
	preserveAspectRatio(query.Get("resize"), size, &options)
	out, err := img.Process(options)
	if !handleUnsupportedFormatError(w, err) {
		return
	}
	ext, mime := getOutputFileMeta(format)
	w.Header().Set("Content-Disposition", "attachment; filename=result"+ext)
	w.Header().Set("Content-Type", mime)
	if _, err := w.Write(out); err != nil {
		log.Println(err)
	}
}
