package deps

import (
	"context"
	"fmt"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var TracerProvider *sdktrace.TracerProvider

func InitTracerProvider(serviceName string) error {
	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))

	if err != nil {
		return fmt.Errorf("создание экспортера трассировки: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		attribute.String("environment", "development"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	TracerProvider = tp

	return nil
}

func Shutdown(ctx context.Context) error {
	if TracerProvider != nil {
		return TracerProvider.Shutdown(ctx)
	}
	return nil
}

func GetTracer(module string) trace.Tracer {
	return otel.Tracer(module)
}

func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := GetTracer("app")
	return tracer.Start(ctx, spanName, opts...)
}
