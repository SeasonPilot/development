package initialization

import (
	"fmt"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"mxshop-api/goods-web/utils/otgrpc"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important;   用户服务发现
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	//grpc lb;
	cfg := global.SrvConfig.ConsulInfo
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", cfg.Host, cfg.Port, global.SrvConfig.GoodsInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// GlobalTracer 默认值是 NoopTracer
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	//fixme: 这里使用 defer conn.Close() 会在 InitSrvConn 本函数执行完后关闭 conn.
	//defer conn.Close()

	global.GoodsClient = proto.NewGoodsClient(conn)
}
