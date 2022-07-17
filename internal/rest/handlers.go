package rest

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rntrp/go-bimg-formpost/internal/config"
)

func handleUnsupportedFormatError(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return true
	case err.Error() == "Unsupported image format":
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType),
			http.StatusUnsupportedMediaType)
		return false
	default:
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return false
	}
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) []byte {
	if !setupFileSizeChecks(w, r) {
		return nil
	}
	memBufSize := coerceMemoryBufferSize(config.GetMemoryBufferSize())
	if err := r.ParseMultipartForm(memBufSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	} else if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	f, fh, err := r.FormFile("img")
	if err != nil {
		http.Error(w, "file 'img' is missing", http.StatusBadRequest)
		return nil
	}
	defer f.Close()
	if fh.Size < minValidFileSize {
		http.Error(w, "file size too small for a valid document",
			http.StatusBadRequest)
		return nil
	}
	input, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return nil
	}
	return input
}

func setupFileSizeChecks(w http.ResponseWriter, r *http.Request) bool {
	clen, err := coerceContentLength(r.Header.Get("Content-Length"))
	if err == nil && clen < minValidFileSize {
		http.Error(w, "http: Content-Length too short for a valid document",
			http.StatusBadRequest)
		return false
	}
	maxReqSize := config.GetMaxRequestSize()
	if maxReqSize >= 0 {
		if err == nil && clen > maxReqSize {
			http.Error(w, "http: Content-Length too large",
				http.StatusRequestEntityTooLarge)
			return false
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxReqSize)
	}
	return true
}
