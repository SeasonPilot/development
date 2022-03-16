package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialization"
	"mxshop-api/user-web/utils/registry/consul"
	myvalidator "mxshop-api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

	// 注册验证器、翻译器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.MobileValidator)
		_ = v.RegisterTranslation("mobile", global.Translator, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

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

	zap.S().Infof("启动服务器，端口：%d", global.SrvConfig.Port)

	//3. 初始化routers
	Routers := initialization.Routers()

	go func() {
		err := Routers.Run(fmt.Sprintf(":%d", global.SrvConfig.Port))
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
