package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuanta7/incidents/pkg/o11y"
)

func main() {
	log.SetOutput(os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pp, err := o11y.NewPrometheusMeterProvider(ctx, "goroutineleak")
	if err != nil {
		panic(err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := pp.Shutdown(shutdownCtx); err != nil {
			log.Printf("meter provider shutdown: %v", err)
		}
	}()

	// pprof at /debug/pprof/goroutine
	// run go tool pprof http://localhost:6060/debug/pprof/goroutine

	http.Handle("/metrics", pp.MetricsHandler())
	http.HandleFunc("/no-leak", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		_, _ = w.Write([]byte("ok"))
	})
	http.HandleFunc("/leak", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan int)
		go func() {
			ch <- 1 // blocks forever, nobody reads, goroutine never exits
		}()
		_, _ = w.Write([]byte("ok but not ok"))
	})

	ln, err := net.Listen("tcp", "localhost:6060")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := &http.Server{}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}()

	go func() {
		log.Printf("server started on %s", ln.Addr())
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")
}
