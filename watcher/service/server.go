package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/ageeknamedslickback/theWatcher/watcher/internal/api/v1"
	"golang.org/x/sync/errgroup"
)

const DefaultTimeout = 3 * time.Second

// NewServer creates a new instance of a http server
func NewServer(addr string) (*http.Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", v1.HealthHandler)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       DefaultTimeout,
		ReadHeaderTimeout: DefaultTimeout,
		WriteTimeout:      DefaultTimeout,
		IdleTimeout:       DefaultTimeout,
	}, nil
}

// Start handles start up and shutdown logic of the service
func Start(f func(ctx context.Context) error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return f(ctx)
	})

	go func() {
		<-ctx.Done()
		slog.Info("forcing shutdown")
		os.Exit(0)
	}()

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("shutting down")
			return
		}
		slog.Error(err.Error())
	}
	os.Exit(0)
}
