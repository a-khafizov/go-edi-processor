package grpc_controllers

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	log *zap.Logger
}

func NewInterceptor(log *zap.Logger) *Interceptor {
	return &Interceptor{log: log}
}

func (i *Interceptor) RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			r := recover()
			if r != nil {
				i.log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func (i *Interceptor) LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		i.log.Info("gRPC method called",
			zap.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
		}

		if err != nil {
			fields = append(fields,
				zap.Error(err),
				zap.Int("code", int(status.Code(err))),
			)
			i.log.Error("gRPC method error", fields...)
		} else {
			i.log.Info("gRPC method success", fields...)
		}

		return resp, err
	}
}
