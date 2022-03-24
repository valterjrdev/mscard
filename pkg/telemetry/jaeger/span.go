package jaeger

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"runtime"
)

func Span(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	return otel.Tracer("").Start(ctx, runtime.FuncForPC(pc).Name())
}
