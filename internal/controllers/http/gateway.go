package http

import (
	"context"
	"net/http"

	"github.com/go-edi-document-processor/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayHandler(ctx context.Context, grpcEndpoint string) (http.Handler, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := proto.RegisterDocumentServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}
