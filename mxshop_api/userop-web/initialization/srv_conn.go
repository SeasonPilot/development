package initialization

import (
	"fmt"

	"mxshop-api/userop-web/global"
	"mxshop-api/userop-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important;   用户服务发现
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	//grpc lb;
	cfg := global.SrvConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.GoodsInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}

	global.GoodsClient = proto.NewGoodsClient(goodsConn)

	// 用户操作服务
	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户操作服务失败】")
	}

	global.UserFavClient = proto.NewUserFavClient(userOpConn)
	global.MessageClient = proto.NewMessageClient(userOpConn)
	global.AddressClient = proto.NewAddressClient(userOpConn)
}
