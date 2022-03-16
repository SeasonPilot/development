package initialization

import (
	"fmt"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/utils/otgrpc"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

func InitTrace() grpc.UnaryServerInterceptor {
	// 集成 tracer
	jaegerConfig := jaegercfg.Configuration{
		ServiceName: global.ServiceConfig.JaegerInfo.Name,
		// 采样
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			// 发送 span 到服务器时是否打印日志
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServiceConfig.JaegerInfo.Host, global.ServiceConfig.JaegerInfo.Port),
		},
	}

	tracer, _, err := jaegerConfig.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	// 过滤掉健康检查，不生成调用链
	opt := otgrpc.IncludingSpans(func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
		return method != "/grpc.health.v1.Health/Check"
	})

	opentracing.SetGlobalTracer(tracer)

	return otgrpc.OpenTracingServerInterceptor(tracer, opt)
}
