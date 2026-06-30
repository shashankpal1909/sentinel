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

	"sentinel/internal/admin"
	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/health"
	"sentinel/internal/logger"
	"sentinel/internal/proxy"
	"sentinel/internal/server"
)

func main() {
	configPath := "example.gateway.yaml"
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		configPath = env
	}

	cfg, err := config.Load(configPath)
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

	mgr, err := app.NewManager(cfg, configPath, l)
	if err != nil {
		slog.Error("Failed to initialize runtime manager", "error", err)
		os.Exit(1)
	}

	slog.Info("Sentinel runtime initialized successfully")
	slog.Info(mgr.GetRuntime().String())

	healthCtx, healthCancel := context.WithCancel(context.Background())
	defer healthCancel()

	checker := health.NewChecker(mgr.GetRuntime(), l)
	checker.Start(healthCtx)
	mgr.SetHealthUpdater(healthCtx, checker)

	p := proxy.New(l)
	srv := server.New(mgr, p, l)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	adminSrv := admin.New(mgr, l)
	adminAddr := fmt.Sprintf("%s:%d", cfg.Admin.Host, cfg.Admin.Port)
	adminHttpServer := &http.Server{
		Addr:    adminAddr,
		Handler: adminSrv,
	}

	// Start Gateway HTTP listener in background goroutine
	go func() {
		slog.Info("Starting Sentinel API Gateway", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Start Admin API HTTP listener in background goroutine
	go func() {
		slog.Info("Starting Sentinel Admin API", "addr", adminAddr)
		if err := adminHttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Admin API server failed", "error", err)
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
		slog.Error("Gateway server shutdown error", "error", err)
	}
	if err := adminHttpServer.Shutdown(ctx); err != nil {
		slog.Error("Admin server shutdown error", "error", err)
	}

	slog.Info("Sentinel API Gateway stopped cleanly")
}
