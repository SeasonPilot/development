package initialization

import (
	"fmt"

	"mxshop-srvs/order_srv/global"
	"mxshop-srvs/order_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important;   用户服务发现
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGoodsSrvConn() {
	//grpc lb;  服务发现
	cfg := global.ServiceConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.ServiceConfig.GoodsInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}

	global.GoodsClient = proto.NewGoodsClient(goodsConn)

	// InitInvSrvConn
	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.ServiceConfig.InvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}

	global.InventoryClient = proto.NewInventoryClient(invConn)
}
