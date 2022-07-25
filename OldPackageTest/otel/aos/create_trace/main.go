package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// 创建Trace&Span（OpenTelemetry Go）
func main() {
	ctx := context.Background()
	// 创建一个tracer
	tracer := otel.Tracer("example/namedtracer/main")

	var span trace.Span
	ctx, span = tracer.Start(ctx, "operation-a")
	defer span.End()

	// 设置Span的Attribute Tag Key Value
	span.SetAttributes(attribute.String("key1", "value1"))

	// 设置Span的Event，可用于记录具体的事件/日
	span.AddEvent("bar", trace.WithAttributes(
		attribute.Bool("key2", true),
		attribute.Int64("key3", 3),
	))
}
