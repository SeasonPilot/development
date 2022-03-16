package middlewares

import (
	"fmt"

	"mxshop-api/goods-web/global"

	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			ServiceName: global.SrvConfig.JaegerInfo.Name,
			// 采样
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				// 发送 span 到服务器时是否打印日志
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.SrvConfig.JaegerInfo.Host, global.SrvConfig.JaegerInfo.Port),
			},
		}

		tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}

		startSpan := tracer.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		c.Set("tracer", tracer)
		c.Set("parentSpan", startSpan)
		c.Next()
	}
}
