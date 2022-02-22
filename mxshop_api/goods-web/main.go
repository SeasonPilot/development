package main

import (
	"fmt"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialization"

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

	err := Routers.Run(fmt.Sprintf(":%d", global.SrvConfig.Port))
	if err != nil {
		zap.S().Panicf("服务启动失败: %s", err.Error())
	}
}
