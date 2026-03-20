package controllers

import (
	"net"

	"github.com/go-edi-document-processor/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcServer struct {
	server *grpc.Server
	log    *logger.Logger
	port   string
}

func NewGrpcServer(log *logger.Logger, port string) *GrpcServer {
	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// Устанавливаем статус SERVING
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	return &GrpcServer{
		server: grpcServer,
		log:    log,
		port:   port,
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	s.log.Zap().Info("gRPC server starting", zap.String("port", s.port))
	return s.server.Serve(lis)
}

func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
	s.log.Zap().Info("gRPC server stopped")
}
