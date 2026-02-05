package graceful

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

type Server interface {
	Run() error
	Shutdown(context.Context) error
}

func RunServer(server Server, timeout time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	notifyCh := make(chan os.Signal, 1)
	defer signal.Stop(notifyCh)

	go func(s Server) {
		if err := s.Run(); err != nil {
			err = fmt.Errorf("error starting REST server: %w", err)
			errCh <- err
		}
	}(server)

	select {
	case err := <-errCh:
		log.Println("Server start error:", err)
		return err
	case <-notifyCh:
		log.Println("Shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, timeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Error during server shutdown:", err)
		return err
	}

	return nil
}
