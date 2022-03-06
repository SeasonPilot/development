package initialization

import (
	"mxshop-api/userop-web/middlewares"
	"mxshop-api/userop-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	// 大写是为了和包名区分开
	Router := gin.Default()

	// 为所有请求添加 Cors 中间件
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/v1")
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)

	return Router
}
