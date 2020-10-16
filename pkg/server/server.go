package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/erumble/go-pd-playground/pkg/logger"
	"golang.org/x/sync/errgroup"
)

// Server provides a function to run an http.HTTPServer with graceful shutdown
type Server interface {
	Serve(ctx context.Context) error
}

type server struct {
	httpServer
	log logger.LeveledLogger

	port            int
	shutdownTimeout time.Duration
}

type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// New does exactly what its name implies
func New(h http.Handler, port int, shutdownTimeout time.Duration, log logger.LeveledLogger) Server {
	httpSrvr := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &server{
		httpServer:      httpSrvr,
		log:             log,
		port:            port,
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *server) Serve(ctx context.Context) error {
	s.log.Infof("Starting server on port: %d", s.port)

	// create a cancel function so we can listen for the done channel
	// this ctx will get overwritten below line, but the cancel
	// function here will close the new context's done channel.
	ctx, cancel := context.WithCancel(context.Background())

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer func() {
		signal.Stop(quit)
		cancel()
	}()

	// Use an errgroup to propagate errors from ListenAndServe
	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		return s.ListenAndServe()
	})

	// Run errg.Wait() in a go routine so we can use a select to listen for the quit channel
	errc := make(chan error)
	go func() {
		defer close(errc)
		errc <- errg.Wait()
	}()

	var errs []error

	// Block until the server exits, or an interrupt is received
	select {
	case <-quit:
		s.log.Debug("Interrupt signal received")
		errs = append(errs, s.stop(ctx))

	// errg.Wait() will call the context's cancel funcion in the event of an error
	case <-ctx.Done():
		s.log.Debug("errg.Wait returned an error")

		// Check for errors from the error group
		for err := range errc {
			errs = append(errs, err)
		}

		errs = append(errs, s.stop(ctx))
	}

	return flattenErrors(errs)
}

func (s *server) stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	s.log.Info("Shutting down server")

	// Ignore context.Cancelled errors
	if err := s.Shutdown(shutdownCtx); err != nil && err != context.Canceled {
		return err
	}
	return nil
}

func flattenErrors(errs []error) error {
	points := []string{}
	for _, err := range errs {
		if err != nil && err != http.ErrServerClosed {
			points = append(points, fmt.Sprintf("* %s", err))
		}
	}

	if len(points) == 0 {
		return nil
	}

	return fmt.Errorf("%d errors occurred:\n\t%s", len(points), strings.Join(points, "\n\t"))
}
