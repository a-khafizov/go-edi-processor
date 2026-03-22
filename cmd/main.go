package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-edi-document-processor/internal/bootstrap/config"
	"github.com/go-edi-document-processor/internal/bootstrap/logger"
	g "github.com/go-edi-document-processor/internal/controllers/grpc_controllers"
	h "github.com/go-edi-document-processor/internal/controllers/http_controllers"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Starting service ...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	log := logger.NewLogger(cfg.LogLevel)
	defer log.Sync()

	startTime := time.Now()

	gatewayCtx := context.Background()
	gatewayHandler, err := h.NewGatewayHandler(gatewayCtx, "localhost:"+cfg.GRPCPort)
	if err != nil {
		log.Error("Failed to create gateway handler",
			zap.Error(err),
		)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: gatewayHandler,
	}

	grpcServer := g.NewGrpcServer(log, cfg.GRPCPort)

	go func() {
		log.Info("Starting HTTP server",
			zap.String("port", cfg.HTTPPort),
		)

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed",
				zap.Error(err),
			)
		}
	}()

	go func() {
		log.Info("Starting gRPC server",
			zap.String("port", cfg.GRPCPort),
		)

		err := grpcServer.Start()
		if err != nil {
			log.Error("gRPC server failed",
				zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down servers...",
		zap.Duration("uptime", time.Since(startTime)),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP server shutdown failed",
			zap.Error(err),
		)
	}
	grpcServer.Stop()

	log.Info("Servers stopped gracefully",
		zap.Duration("total_uptime", time.Since(startTime)),
	)
}
