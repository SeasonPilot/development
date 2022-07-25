package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

// 配置接入点完成数据上报（OpenTelemetry Go）
func installExportPipeline(ctx context.Context) func() {

	// 1.创建并定义资源
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// 请设置对应的环境变量并将YOUR SERVICE NAME，YOUR SERVICE NAMESPACE，
			// YOUR SERVICE DEPLOY ENVIRONMENT设置成对应的服务名，命名空间，环境。默认Default
			// 有相同的名称，命名空间，及环境的观测对象被认为是同一个服务，运维中心会自动对接入服务进行服务名的拼接，格式为{服务名}.{命名空间}.{环境}
			semconv.ServiceNameKey.String("YOUR SERVICE NAME"),
			semconv.ServiceNamespaceKey.String("YOUR SERVICE NAMESPACE"),
			semconv.DeploymentEnvironmentKey.String("YOUR SERVICE DEPLOY ENVIRONMENT"),
		),
	)
	handleErr(err, "failed to create resource")

	// 2.请将AOS_COLLECTOR_ENDPOINT设置为从
	otelAgentAddr, ok := os.LookupEnv("AOS_COLLECTOR_ENDPOINT")
	if !ok {
		otelAgentAddr = "0.0.0.0:4317"
	}

	// 3.设置一个新的GRPC Trace Client并绑定之前设置的上报端口
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	traceExp, err := otlptrace.New(ctx, traceClient)
	handleErr(err, "Failed to create the collector trace exporter")

	// 4. 设置批量上报，设置客户端采样率，并绑定之前定义的资源
	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// 5. 设置全局Propagator
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func main() {
	ctx := context.Background()
	// Registers a tracer Provider globally.
	cleanup := installExportPipeline(ctx)
	defer cleanup()
}
