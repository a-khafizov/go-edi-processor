package deps

import (
	"net"

	proto "github.com/go-edi-document-processor/api/proto/gen"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	server                     *grpc.Server
	port                       string
	logger                     *zap.Logger
	protoDocumentServiceServer proto.DocumentServiceServer
}

func NewGrpcServer(logger *zap.Logger, port string, tracer trace.Tracer, protoDocumentServiceServer proto.DocumentServiceServer) *GrpcServer {
	interceptor := NewInterceptor(logger, tracer)

	recovery := interceptor.RecoveryInterceptor()
	tracing := interceptor.TracingInterceptor()
	logging := interceptor.LoggingInterceptor()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery,
			tracing,
			logging,
		),
	)

	proto.RegisterDocumentServiceServer(grpcServer, protoDocumentServiceServer)

	return &GrpcServer{
		server:                     grpcServer,
		port:                       port,
		logger:                     logger,
		protoDocumentServiceServer: protoDocumentServiceServer,
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
