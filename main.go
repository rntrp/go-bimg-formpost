package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/rntrp/go-bimg-formpost/internal/config"
	"github.com/rntrp/go-bimg-formpost/internal/rest"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	slog.Info("Loading GO-BIMG-FORMPOST...")
	config.Load()
}

func main() {
	if err := start(); err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("Bye.")
	}
}

func start() error {
	srvout := make(chan error, 1)
	signal := make(chan os.Signal, 1)
	srv := server(signal)
	go shutdownMonitor(signal, srvout, srv)
	slog.Info("Starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-srvout
}

func server(sig chan os.Signal) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("GET /", rest.Index)
	r.HandleFunc("GET /index.html", rest.Index)
	r.HandleFunc("GET /live", rest.Live)
	r.HandleFunc("POST /convert", rest.Convert)
	if config.IsEnablePrometheus() {
		r.Handle("/metrics", promhttp.Handler())
	}
	if config.IsEnableShutdown() {
		r.HandleFunc("POST /shutdown", shutdownFn(sig))
	}
	h := httplog.Handler(httplog.NewLogger("GO-BIMG-FORMPOST", httplog.Options{
		Concise:         true,
		JSON:            false,
		RequestHeaders:  false,
		TimeFieldFormat: time.RFC3339,
	}))(r)
	return &http.Server{Addr: config.GetTCPAddress(), Handler: h}
}

func shutdownFn(sig chan os.Signal) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		go func() { sig <- os.Interrupt }()
	}
}

func shutdownMonitor(sig chan os.Signal, out chan error, srv *http.Server) {
	timeout := config.GetShutdownTimeout()
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sigName := (<-sig).String()
	slog.Info("Shutdown signal received", "signal", sigName)
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	out <- srv.Shutdown(ctx)
}
