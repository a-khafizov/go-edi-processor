package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-edi-document-processor/internal/config"
	"github.com/go-edi-document-processor/internal/controllers"
	"github.com/go-edi-document-processor/internal/logger"
	"github.com/go-edi-document-processor/internal/middleware"
	"go.uber.org/zap"
)

func main() {
	// Инициализация конфигурации
	log.Printf("Starting service with environment: %s", config.Environment())

	// Инициализация логгера
	err := logger.InitGlobal(config.LogLevel(), config.IsDevelopment())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger := logger.GetGlobal()
	defer logger.Sync()

	// Создание HTTP сервера
	gin.SetMode(gin.ReleaseMode)
	if config.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	m := middleware.NewMiddleware(logger)
	r.Use(m.Recovery())
	r.Use(m.RequestLogger())

	restController := controllers.NewRestController(logger, m)
	restController.RegisterRoutes(r)

	httpPort := config.HTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}
	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: r,
	}

	// Создание gRPC сервера
	grpcPort := config.GRPCPort()
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcServer := controllers.NewGrpcServer(logger, grpcPort)

	// Запуск серверов в горутинах
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

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Zap().Info("Shutting down servers...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Zap().Error("HTTP server shutdown error", zap.Error(err))
	}
	grpcServer.Stop()

	logger.Zap().Info("Servers stopped gracefully")
}
