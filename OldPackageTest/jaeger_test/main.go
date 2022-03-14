package main

import (
	"time"

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

	span := tracer.StartSpan("go-grpc-web")
	<-time.After(time.Second)
	defer span.Finish()
}
