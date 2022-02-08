package main

import (
	"fmt"

	"development/mxshop_api/user-web/initialization"

	"go.uber.org/zap"
)

func main() {
	port := 8021
	//1. 初始化logger
	initialization.InitLogger()

	zap.S().Infof("启动服务器，端口：%d", port)

	//3. 初始化routers
	Routers := initialization.Routers()

	err := Routers.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panicf("服务启动失败: %s", err.Error())
	}
}
