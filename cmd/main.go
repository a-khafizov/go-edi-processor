package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	g "github.com/go-edi-document-processor/internal/controllers/grpc"
	h "github.com/go-edi-document-processor/internal/controllers/http"
	"github.com/go-edi-document-processor/internal/infrastructure/config"
	"github.com/go-edi-document-processor/internal/infrastructure/logger"
	"github.com/go-edi-document-processor/internal/infrastructure/tracing"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Starting service with environment: %s", config.Environment())

	err := logger.InitGlobal(config.LogLevel(), config.IsDevelopment())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger := logger.GetGlobal()
	defer logger.Sync()

	shutdownTracing, err := tracing.InitTracing("go-edi-document-processor", io.Discard)
	if err != nil {
		logger.Zap().Fatal("Failed to initialize tracing", zap.Error(err))
	}
	defer shutdownTracing()

	httpPort := config.HTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}
	grpcPort := config.GRPCPort()
	if grpcPort == "" {
		grpcPort = "50051"
	}

	gatewayCtx := context.Background()
	gatewayHandler, err := h.NewGatewayHandler(gatewayCtx, "localhost:"+grpcPort)
	if err != nil {
		logger.Zap().Fatal("Failed to create gateway handler", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: gatewayHandler,
	}

	grpcServer := g.NewGrpcServer(logger, grpcPort)

	go func() {
		logger.Zap().Info("Starting HTTP server", zap.String("port", httpPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Zap().Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Zap().Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Zap().Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Zap().Error("HTTP server shutdown error", zap.Error(err))
	}
	grpcServer.Stop()

	logger.Zap().Info("Servers stopped gracefully")
}
