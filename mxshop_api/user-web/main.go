package main

import (
	"fmt"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialization"
	myvalidator "mxshop-api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	//1. 初始化logger
	initialization.InitLogger()
	initialization.InitConfig()
	initialization.InitTrans("zh")
	initialization.InitSrvConn()

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

	zap.S().Infof("启动服务器，端口：%d", global.SrvConfig.Port)

	//3. 初始化routers
	Routers := initialization.Routers()

	err := Routers.Run(fmt.Sprintf(":%d", global.SrvConfig.Port))
	if err != nil {
		zap.S().Panicf("服务启动失败: %s", err.Error())
	}
}
