package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"golang.org/x/sync/errgroup"
)

const (
	defaultAppPort        = "8080"
	httpReadTimeout       = 10 * time.Second
	httpReadHeaderTimeout = 3 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpMaxHeaderBytes    = 1 << 20
	shutdownTimeout       = 10 * time.Second
)

func main() {
	logger := newLogger()
	slog.SetDefault(logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, logger); err != nil {
		logger.Error("application terminated unexpectedly", "error", err)
		os.Exit(1)
	}
	logger.Info("application exited successfully")
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
}

func run(ctx context.Context, logger *slog.Logger) error {
	router := newRouter(logger)
	server := newServer(router, logger)
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server startup failed: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		logger.Info("server shutdown complete")
		return nil
	})
	return g.Wait()
}

func newRouter(logger *slog.Logger) *gin.Engine {
	if os.Getenv(gin.EnvGinMode) == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent: true,
		WithRequestID: true,
		Filters: []sloggin.Filter{
			sloggin.IgnorePath("/health"),
		},
	}))
	router.Use(gin.Recovery())
	router.GET("/health", healthCheckHandler)
	return router
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func newServer(handler http.Handler, logger *slog.Logger) *http.Server {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = defaultAppPort
	}
	return &http.Server{
		Addr:              net.JoinHostPort("", appPort),
		Handler:           handler,
		ReadTimeout:       httpReadTimeout,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
		MaxHeaderBytes:    httpMaxHeaderBytes,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
