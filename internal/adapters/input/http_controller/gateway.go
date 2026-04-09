package gateway

import (
	"context"
	"net/http"

	proto "github.com/go-edi-document-processor/api/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewHttpControllers(ctx context.Context, grpcEndpoint string) (http.Handler, error) {
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
