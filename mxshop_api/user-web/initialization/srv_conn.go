package initialization

import (
	"fmt"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important;   用户服务发现
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	cfg := global.SrvConfig.ConsulInfo
	//grpc lb;  consul 用户服务发现
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.UserInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	//fixme: 这里使用 defer conn.Close() 会在 InitSrvConn 本函数执行完后关闭 conn.
	//defer conn.Close()

	global.UserClient = proto.NewUserClient(conn)
}

// InitSrvConn2 第一个版本,存档
func InitSrvConn2() {
	// 创建 consul Client
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulInfo.Host, global.SrvConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return
	}

	// 服务发现
	result, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.SrvConfig.Name))
	if err != nil {
		panic(err)
		return
	}

	if len(result) == 0 {
		panic(fmt.Errorf("服务发现错误，没有该服务"))
		return
	}

	var (
		addr string
		port int
	)

	for _, service := range result {
		addr = service.Address
		port = service.Port
	}

	if addr == "" || port == 0 {
		zap.S().Errorf("addr: %s ,  port:%d", addr, port)
		panic(fmt.Errorf("服务发现错误，没有该服务"))
		return
	}

	// non-blocking dial
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicf("连接用户服务失败 %s", err.Error())
		return
	}
	global.UserClient = proto.NewUserClient(conn)
}
