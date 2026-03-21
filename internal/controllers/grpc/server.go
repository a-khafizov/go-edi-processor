package grpc

import (
	"net"

	"github.com/go-edi-document-processor/api/proto"
	"github.com/go-edi-document-processor/internal/infrastructure/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server *grpc.Server
	log    *logger.Logger
	port   string
}

func NewGrpcServer(log *logger.Logger, port string) *GrpcServer {
	otelHandler := otelgrpc.NewServerHandler()
	recovery := RecoveryInterceptor(log)
	logging := LoggingInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelHandler),
		grpc.ChainUnaryInterceptor(
			recovery,
			logging,
		),
	)

	docService := NewDocumentService(log)
	proto.RegisterDocumentServiceServer(grpcServer, docService)

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
