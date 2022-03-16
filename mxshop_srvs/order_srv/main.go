package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-srvs/order_srv/global"
	"mxshop-srvs/order_srv/handler"
	"mxshop-srvs/order_srv/initialization"
	"mxshop-srvs/order_srv/proto"
	"mxshop-srvs/order_srv/utils/otgrpc"
	"mxshop-srvs/order_srv/utils/registry/consul"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitDB()
	initialization.InitRedisClient()
	initialization.InitRedSync()
	initialization.InitGoodsSrvConn()

	//freePort, err := utils.GetFreePort()
	//if err != nil {
	//	panic(err)
	//}
	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	zap.S().Infof("ip: %s, port: %d", *ip, *port)

	// 集成 tracer
	cfg := jaegercfg.Configuration{
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

	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)

	// 过滤掉健康检查，不生成调用链
	opt := otgrpc.IncludingSpans(func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
		return method != "/grpc.health.v1.Health/Check"
	})
	// 注册用户服务; 集成 OpenTracingServerInterceptor
	g := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer, opt)))
	proto.RegisterOrderServer(g, &handler.OrderServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}

	// 注册 grpc 服务健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	go func() {
		// 启动服务
		err = g.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	// 服务注册
	var rc consul.RegisterClient
	srvID := uuid.New().String()
	rc = consul.NewConsulClient(global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)
	err = rc.Register(srvID,
		global.ServiceConfig.Name,
		global.ServiceConfig.Tags,
		*port,
		global.ServiceConfig.Host,
	)
	if err != nil {
		panic(err)
	}

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(primitive.NamesrvAddr{"172.19.30.30:9876"}),
		// 通过 GroupName 可以达到负载均衡的效果
		consumer.WithGroupName("mxshop-order"),
	)

	if err = c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Println("获得消息失败")
		return
	}

	_ = c.Start()

	// 优雅退出; deregister 服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = c.Shutdown()

	err = rc.Deregister(srvID)
	if err != nil {
		return
	}
}
