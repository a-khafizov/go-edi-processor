package grpc

import (
	"context"
	"time"

	"github.com/go-edi-document-processor/internal/infrastructure/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	spanContext := span.SpanContext()
	if spanContext.HasTraceID() {
		return spanContext.TraceID().String()
	}
	return ""
}

func getSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	spanContext := span.SpanContext()
	if spanContext.HasSpanID() {
		return spanContext.SpanID().String()
	}
	return ""
}

func RecoveryInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				traceID := getTraceID(ctx)
				spanID := getSpanID(ctx)
				log.Zap().Error("panic recovered",
					zap.Any("panic", r),
					zap.String("trace_id", traceID),
					zap.String("span_id", spanID),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func LoggingInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		traceID := getTraceID(ctx)
		spanID := getSpanID(ctx)
		start := time.Now()
		log.Zap().Info("gRPC method called",
			zap.String("method", info.FullMethod),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		level := log.Zap().Info
		if err != nil {
			level = log.Zap().Error
		}
		level("gRPC method completed",
			zap.String("method", info.FullMethod),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return resp, err
	}
}
