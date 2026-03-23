package deps

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	log    *zap.Logger
	tracer trace.Tracer
}

func NewInterceptor(log *zap.Logger, tracer trace.Tracer) *Interceptor {
	return &Interceptor{log: log, tracer: tracer}
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

		span := trace.SpanFromContext(ctx)
		var traceID string
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		} else {
			traceID = "no-trace"
		}

		i.log.Info("gRPC method called",
			zap.String("method", info.FullMethod),
			zap.String("trace_id", traceID),
		)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("trace_id", traceID),
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

func (i *Interceptor) TracingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if i.tracer == nil {
			return handler(ctx, req)
		}

		spanName := info.FullMethod
		ctx, span := i.tracer.Start(ctx, spanName)
		defer span.End()

		span.SetAttributes(
			attribute.String("rpc.method", info.FullMethod),
			attribute.String("rpc.system", "grpc"),
		)

		return handler(ctx, req)
	}
}
