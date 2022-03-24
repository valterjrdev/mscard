package jaeger

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"os"
)

func Setup(name string, namespace string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return err
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceNamespaceKey.String(namespace),
			semconv.ServiceInstanceIDKey.String(hostname),
		)),
	)

	otel.SetTracerProvider(provider)
	return nil
}
