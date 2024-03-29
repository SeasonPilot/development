package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		ServiceName: "mxshop",
		// 采样
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			// 发送 span 到服务器时是否打印日志
			LogSpans:           true,
			LocalAgentHostPort: "172.19.30.31:6831",
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// 设置全局 Tracer  核心代码从这开始到最后
	opentracing.SetGlobalTracer(tracer)

	parentSpan := opentracing.StartSpan("main") // 直接调用 opentracing 包的 StartSpan，不用 cfg.NewTracer 返回的 tracer 了

	span := opentracing.StartSpan("FunA", opentracing.ChildOf(parentSpan.Context()))
	<-time.After(time.Millisecond * 500)
	span.Finish()

	<-time.After(time.Millisecond * 100)

	span2 := opentracing.StartSpan("FunB", opentracing.ChildOf(span.Context()))
	<-time.After(time.Second)
	span2.Finish()

	parentSpan.Finish()
}
