package tracer

import (
	"context"
	"fmt"
	"time"

	gcp "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func Must(f func() error, err error) func() error {
	if err != nil {
		panic(err)
	}

	return f
}

func Setup(projectID, serviceName, revision string, timeout time.Duration) (func() error, error) {
	exporter, err := gcp.New(gcp.WithProjectID(projectID))
	if err != nil {
		return nil, fmt.Errorf("new exporter: %v", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(revision),
			),
		),
	)

	otel.SetTracerProvider(provider)

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := provider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("provider force flush: %v", err)
		}

		if err := provider.Shutdown(ctx); err != nil {
			return fmt.Errorf("provider shutdown: %v", err)
		}

		return nil
	}, nil
}

func Context(ctx context.Context, traceID, spanID string, traceTrue bool) (context.Context, error) {
	tID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, fmt.Errorf("traceID from hex(%v): %v", traceID, err)
	}

	sID, err := trace.SpanIDFromHex(spanID)
	if err != nil {
		return nil, fmt.Errorf("spanID from hex(%v): %v", spanID, err)
	}

	flags := trace.TraceFlags(00)
	if traceTrue {
		flags = 01
	}

	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: flags,
		Remote:     false,
	})), nil
}
