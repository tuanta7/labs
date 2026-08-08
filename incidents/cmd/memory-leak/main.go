package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tuanta7/incidents/pkg/o11y"
	"go.uber.org/zap"
)

var (
	mu     sync.Mutex
	leaked [][]byte
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := o11y.NewLogger(ctx, "memoryleak")
	if err != nil {
		panic(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := logger.Shutdown(shutdownCtx); err != nil {
			logger.Error("logger shutdown", zap.Error(err))
		}
	}()

	pp, err := o11y.NewPrometheusMeterProvider(ctx, "memoryleak")
	if err != nil {
		panic(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := pp.Shutdown(shutdownCtx); err != nil {
			logger.Error("meter provider shutdown", zap.Error(err))
		}
	}()

	// pprof at /debug/pprof/heap
	// run go tool pprof http://localhost:6061/debug/pprof/heap

	http.Handle("/metrics", pp.MetricsHandler())
	http.HandleFunc("/no-leak", func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 1<<20) // 1MB, discarded after the request, GC'd
		_ = b
		_, _ = w.Write([]byte("ok"))
	})
	http.HandleFunc("/leak", func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 1<<20) // 1MB, held forever, never freed

		mu.Lock()
		leaked = append(leaked, b)
		total := len(leaked)
		mu.Unlock()

		logger.Info("leaked chunk", zap.Int("total_chunks", total))
		_, _ = w.Write([]byte("ok but not ok"))
	})

	ln, err := net.Listen("tcp", "localhost:6061")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	srv := &http.Server{}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("server started", zap.String("addr", ln.Addr().String()))
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
}
