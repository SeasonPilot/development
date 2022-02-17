package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-srvs/user_srv/global"
	"mxshop-srvs/user_srv/handler"
	"mxshop-srvs/user_srv/initialization"
	"mxshop-srvs/user_srv/proto"
	"mxshop-srvs/user_srv/utils"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitDB()

	freePort, err := utils.GetFreePort()
	if err != nil {
		panic(err)
		return
	}

	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", freePort, "端口号")

	flag.Parse()
	zap.S().Info(*ip, *port)

	// 注册用户服务
	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handler.UserServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}

	// 注册 grpc 服务健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	// 注册服务到 consul, 即服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	srvID := uuid.New().String()
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      srvID,
		Name:    global.ServiceConfig.Name,
		Tags:    []string{"primary", "v1"},
		Port:    *port,
		Address: "10.94.62.100",
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d", "10.94.62.100", *port),
			DeregisterCriticalServiceAfter: "15s",
		},
	})
	if err != nil {
		panic(err)
	}

	go func() {
		// 启动服务
		err = g.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	// 优雅退出; deregister 服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = client.Agent().ServiceDeregister(srvID)
	if err != nil {
		zap.S().Errorf("注销服务失败:%s %s", global.ServiceConfig.Name, srvID)
		return
	}
	zap.S().Info("注销成功")
}
