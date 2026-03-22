package grpc_controllers

import (
	"net"

	"github.com/go-edi-document-processor/api/proto"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	server *grpc.Server
	port   string
	logger *zap.Logger
}

func NewGrpcServer(logger *zap.Logger, port string) *GrpcServer {
	interceptor := NewInterceptor(logger)

	recovery := interceptor.RecoveryInterceptor()
	logging := interceptor.LoggingInterceptor()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery,
			logging,
		),
	)

	docService := NewDocumentService(logger)
	proto.RegisterDocumentServiceServer(grpcServer, docService)

	return &GrpcServer{
		server: grpcServer,
		port:   port,
		logger: logger,
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	return s.server.Serve(lis)
}

func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
}
