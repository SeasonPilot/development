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

	parentSpan := tracer.StartSpan("main")

	span := tracer.StartSpan("FunA", opentracing.ChildOf(parentSpan.Context()))
	<-time.After(time.Millisecond * 500)
	span.Finish()

	<-time.After(time.Millisecond * 100)

	span2 := tracer.StartSpan("FunB", opentracing.ChildOf(span.Context()))
	<-time.After(time.Second)
	span2.Finish()

	parentSpan.Finish()
}
