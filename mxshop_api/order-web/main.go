package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mxshop-api/order-web/global"
	"mxshop-api/order-web/initialization"
	"mxshop-api/order-web/registry/consul"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	//1. 初始化logger
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitTrans("zh")
	initialization.InitSrvConn()

	//if !initialization.GetEnvInfo("MXSHOP_DEBUG") {
	//	port, err := utils.GetFreePort()
	//	if err != nil {
	//		panic(err)
	//		return
	//	}
	//	global.SrvConfig.Port = port
	//}

	zap.S().Infof("启动服务器，端口：%d", global.SrvConfig.Port)

	//3. 初始化routers
	Routers := initialization.Routers()

	// 服务注册
	var rc consul.RegisterClient
	srvID := uuid.New().String()
	rc = consul.NewConsulClient(global.SrvConfig.ConsulInfo.Host, global.SrvConfig.ConsulInfo.Port)
	err := rc.Register(srvID,
		global.SrvConfig.Name,
		global.SrvConfig.Tags,
		global.SrvConfig.Port,
		global.SrvConfig.Address,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		err = Routers.Run(fmt.Sprintf(":%d", global.SrvConfig.Port))
		if err != nil {
			zap.S().Panicf("服务启动失败: %s", err.Error())
		}
	}()

	// 优雅退出; deregister 服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = rc.Deregister(srvID)
	if err != nil {
		return
	}
}
