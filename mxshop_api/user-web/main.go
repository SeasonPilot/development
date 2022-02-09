package main

import (
	"fmt"

	"development/mxshop_api/user-web/global"
	"development/mxshop_api/user-web/initialization"

	"go.uber.org/zap"
)

func main() {
	//1. 初始化logger
	initialization.InitLogger()
	initialization.InitConfig()

	zap.S().Infof("启动服务器，端口：%d", global.SrvConfig.Port)

	//3. 初始化routers
	Routers := initialization.Routers()

	err := Routers.Run(fmt.Sprintf(":%d", global.SrvConfig.Port))
	if err != nil {
		zap.S().Panicf("服务启动失败: %s", err.Error())
	}
}
