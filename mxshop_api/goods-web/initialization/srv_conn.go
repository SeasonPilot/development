package initialization

import (
	"fmt"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important;   用户服务发现
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	//grpc lb;
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d?", global.SrvConfig.GoodsInfo.Host, global.SrvConfig.GoodsInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	//fixme: 这里使用 defer conn.Close() 会在 InitSrvConn 本函数执行完后关闭 conn.
	//defer conn.Close()

	global.GoodsClient = proto.NewGoodsClient(conn)
}
