package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/ageeknamedslickback/theWatcher/watcher/service"
	"golang.org/x/sync/errgroup"
)

func main() {
	service.Start(func(ctx context.Context) error {
		port, ok := os.LookupEnv("PORT")
		if !ok {
			slog.Error("PORT env var is unset")
			return errors.New("PORT env var is required")
		} else if _, err := strconv.Atoi(port); err != nil {
			slog.Error(port + " is not a valid port number")
			return err
		}

		addr := fmt.Sprintf(":%s", port)

		server, err := service.NewServer(addr)
		if err != nil {
			return err
		}

		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			slog.Info(fmt.Sprintf("Starting server on %s", addr))
			if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		})
		g.Go(func() error {
			<-ctx.Done()
			return server.Shutdown(ctx)
		})

		return g.Wait()
	})
}
