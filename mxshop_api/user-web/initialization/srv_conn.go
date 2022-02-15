package initialization

import (
	"fmt"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
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
