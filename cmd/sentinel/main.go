package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/health"
	"sentinel/internal/logger"
	"sentinel/internal/proxy"
	"sentinel/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if err := config.Validate(cfg); err != nil {
		slog.Error("Invalid config", "error", err)
		os.Exit(1)
	}

	l := logger.Init("info")
	logger.PrintBanner()

	rt, err := app.Build(cfg)
	if err != nil {
		slog.Error("Failed to build runtime", "error", err)
		os.Exit(1)
	}

	slog.Info("Sentinel runtime initialized successfully")
	slog.Info(rt.String())

	healthCtx, healthCancel := context.WithCancel(context.Background())
	defer healthCancel()

	checker := health.NewChecker(rt, l)
	checker.Start(healthCtx)

	p := proxy.New(l)
	srv := server.New(rt, p, l)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	// Start HTTP listener in background goroutine to allow signal interception
	go func() {
		slog.Info("Starting Sentinel API Gateway", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Intercept termination signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("Shutting down Sentinel API Gateway gracefully...")

	healthCancel()
	checker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("Sentinel API Gateway stopped cleanly")
}
