package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-srvs/userop_srv/global"
	"mxshop-srvs/userop_srv/handler"
	"mxshop-srvs/userop_srv/initialization"
	"mxshop-srvs/userop_srv/proto"
	"mxshop-srvs/userop_srv/utils/registry/consul"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitDB()

	//freePort, err := utils.GetFreePort()
	//if err != nil {
	//	panic(err)
	//}
	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	zap.S().Infof("ip: %s, port: %d", *ip, *port)

	// 注册服务
	g := grpc.NewServer()
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

	// 优雅退出; deregister 服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = rc.Deregister(srvID)
	if err != nil {
		return
	}
}
