package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-srvs/inventory_srv/global"
	"mxshop-srvs/inventory_srv/handler"
	"mxshop-srvs/inventory_srv/initialization"
	"mxshop-srvs/inventory_srv/proto"
	"mxshop-srvs/inventory_srv/utils/registry/consul"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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
	initialization.InitRedisClient()
	initialization.InitRedSync()

	//freePort, err := utils.GetFreePort()
	//if err != nil {
	//	panic(err)
	//}
	ip := flag.String("ip", "0.0.0.0", "ip 地址")
	port := flag.Int("port", 50053, "端口号")

	flag.Parse()
	zap.S().Infof("ip: %s, port: %d", *ip, *port)

	// 注册用户服务
	g := grpc.NewServer()
	proto.RegisterInventoryServer(g, &handler.InventoryServer{})

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

	// 监听库存归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(primitive.NamesrvAddr{"172.19.30.30:9876"}),
		// 通过 GroupName 可以达到负载均衡的效果
		consumer.WithGroupName("mxshop-inventory"),
	)

	if err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("获得消息失败")
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
