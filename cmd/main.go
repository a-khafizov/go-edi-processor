package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-edi-document-processor/internal/adapters/input/grpc_controller"
	gateway "github.com/go-edi-document-processor/internal/adapters/input/http_controller"

	adapters "github.com/go-edi-document-processor/internal/adapters/output"
	"github.com/go-edi-document-processor/internal/core/services"
	"github.com/go-edi-document-processor/internal/deps"
	"github.com/oagudo/outbox"
	"go.uber.org/zap"
)

func main() {
	log.Printf("Starting service ...")

	cfg, err := deps.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	logger := deps.InitLogger(cfg.LogLevel)
	defer logger.Sync()

	serviceName := "go-edi-processor"
	if err := deps.InitTracerProvider(serviceName); err != nil {
		logger.Error("Failed to initialize tracing",
			zap.Error(err),
		)
		os.Exit(1)
	}
	defer deps.Shutdown(context.Background())

	db, err := deps.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	redisClient, err := deps.InitRedis(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	startTime := time.Now()

	docRepository := adapters.NewDocumentRepository(db)

	cacheRepository := adapters.NewRedisCache(redisClient)

	dbCtx := outbox.NewDBContext(db, outbox.SQLDialectPostgres)

	outboxService, err := adapters.NewOutboxService(db, dbCtx, docRepository)
	if err != nil {
		logger.Fatal("Failed to create outbox service", zap.Error(err))
	}

	kafkaPublisher, err := adapters.NewKafkaPublisher(cfg)
	if err != nil {
		logger.Fatal("Failed to create Kafka publisher", zap.Error(err))
	}
	defer kafkaPublisher.Close()

	outboxReader := adapters.NewOutboxReader(dbCtx, kafkaPublisher, logger)
	outboxReader.Start()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := outboxReader.Stop(ctx); err != nil {
			logger.Error("Failed to stop outbox reader", zap.Error(err))
		}
	}()

	kafkaConsumer, err := adapters.NewKafkaConsumer(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}
	defer kafkaConsumer.Close()

	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()
	go kafkaConsumer.Start(consumerCtx)

	docService := services.NewDocumentService(docRepository, outboxService, cacheRepository)

	protoDocumentServiceServer := grpc_controller.NewProtoDocumentServiceServer(docService)

	gatewayCtx := context.Background()
	gateway, err := gateway.NewHttpControllers(gatewayCtx, "localhost:"+cfg.GRPCPort)
	if err != nil {
		logger.Error("Failed to create gateway handler",
			zap.Error(err),
		)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: gateway,
	}

	tracer := deps.GetTracer("grpc")
	grpcServer := deps.NewGrpcServer(logger, cfg.GRPCPort, tracer, protoDocumentServiceServer)

	go func() {
		logger.Info("Starting HTTP server",
			zap.String("port", cfg.HTTPPort),
		)

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed",
				zap.Error(err),
			)
		}
	}()

	go func() {
		logger.Info("Starting gRPC server",
			zap.String("port", cfg.GRPCPort),
		)

		err := grpcServer.Start()
		if err != nil {
			logger.Error("gRPC server failed",
				zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...",
		zap.Duration("uptime", time.Since(startTime)),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown failed",
			zap.Error(err),
		)
	}
	grpcServer.Stop()

	logger.Info("Servers stopped gracefully",
		zap.Duration("total_uptime", time.Since(startTime)),
	)
}
