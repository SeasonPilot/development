package initialization

import (
	"fmt"

	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"

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

	// 库存服务
	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.InvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}

	global.InvClient = proto.NewInventoryClient(invConn)

	//  订单服务
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】")
	}

	global.OrderClient = proto.NewOrderClient(orderConn)
}
