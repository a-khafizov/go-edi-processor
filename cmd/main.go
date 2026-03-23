package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adapters_grpc "github.com/go-edi-document-processor/internal/adapters/primary/grpc_controllers"
	adapters_http "github.com/go-edi-document-processor/internal/adapters/primary/http_controllers"

	adapters "github.com/go-edi-document-processor/internal/adapters/secondary"
	"github.com/go-edi-document-processor/internal/core/services"
	"github.com/go-edi-document-processor/internal/deps"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Starting service ...")

	cfg, err := deps.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	log := deps.NewLogger(cfg.LogLevel)
	defer log.Sync()

	serviceName := "go-edi-document-processor"
	if err := deps.InitTracerProvider(serviceName); err != nil {
		log.Error("Failed to initialize tracing",
			zap.Error(err),
		)
		os.Exit(1)
	}
	defer deps.Shutdown(context.Background())

	startTime := time.Now()

	docRepository := adapters.NewDocumentRepository()
	// outboxRepository := adapters.NewOutboxRepository()

	docService := services.NewDocumentService(docRepository)
	// outboxService := services.NewOutboxService(outboxRepository)

	protoDocumentServiceServer := adapters_grpc.NewProtoDocumentServiceServer(docService)

	gatewayCtx := context.Background()
	httpController, err := adapters_http.NewHttpControllers(gatewayCtx, "localhost:"+cfg.GRPCPort)
	if err != nil {
		log.Error("Failed to create gateway handler",
			zap.Error(err),
		)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: httpController,
	}

	tracer := deps.GetTracer("grpc")
	grpcServer := deps.NewGrpcServer(log, cfg.GRPCPort, tracer, protoDocumentServiceServer)

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
