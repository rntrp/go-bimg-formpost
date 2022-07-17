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
	for _, err := range [...]error{errW, errH, errF} {
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
	}
	src := handleFileUpload(w, r)
	if src == nil {
		return
	}
	options := bimg.Options{
		Width:         width,
		Height:        height,
		Type:          format,
		Compression:   9,
		Quality:       99,
		Extend:        bimg.ExtendWhite,
		Background:    bimg.Color{R: 0xFF, G: 0xFF, B: 0xFF},
		Embed:         true,
		Enlarge:       true,
		StripMetadata: true,
		NoProfile:     true,
	}
	out, err := bimg.NewImage(src).Process(options)
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
