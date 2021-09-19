package main

import (
	"log"
	"net/http"

	"github.com/rntrp/bimg-rest/internal/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/h2non/bimg"
)

const address = ":8080"

func main() {
	log.Printf("libvips version: %s", bimg.VipsVersion)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", rest.Welcome)
	r.Post("/scale", rest.Scale)
	http.ListenAndServe(address, r)
}
