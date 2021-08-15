package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
)

const maxMemory int64 = 1024 * 1024 * 64
const maxFileSize int64 = 1024 * 1024 * 256

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func Scale(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	format := getFormat(query.Get("format"))
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	f, fh, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File 'image' could not be processed."))
		return
	}
	defer f.Close()
	if fh.Size > maxFileSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("Max file size is 256 MiB."))
		return
	} 
	printMemUsage(1)
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(buf) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File 'image' is empty."))
		return
	}
	printMemUsage(2)
	options := getOptions(width, height, format)
	o, err := bimg.NewImage(buf).Process(options)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	printMemUsage(3)
	w.Write(o)
}

func getFormat(format string) bimg.ImageType {
	switch (strings.ToLower(format)) {
	case "png": return bimg.PNG
	case "webp": return bimg.WEBP
	case "gif": return bimg.GIF
	case "heif", "heic": return bimg.HEIF
	case "avif": return bimg.AVIF
	default: return bimg.JPEG
	}
}

func getOptions(width, height int, format bimg.ImageType) bimg.Options {
	return bimg.Options {
		Width:         width,
		Height:        height,
		Type:          format,
		Compression:   9,
		Quality:       99,
		Extend:        bimg.ExtendWhite,
		Background:    bimg.Color{R:0xFF, G:0xFF, B:0xFF},
		Embed:         true,
		Enlarge:       true,
		StripMetadata: true,
		NoProfile:     true,
	}
}

func printMemUsage(step int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Step %d: Alloc = %v MiB", step, bToMiB(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMiB(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMiB(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Println()
}

func bToMiB(b uint64) string {
	f := float64(b) / 1_048_576
	return strconv.FormatFloat(f, 'f', 3, 64)
}
