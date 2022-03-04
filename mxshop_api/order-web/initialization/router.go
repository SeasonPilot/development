package initialization

import (
	"mxshop-api/order-web/middlewares"
	"mxshop-api/order-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	// 大写是为了和包名区分开
	Router := gin.Default()

	// 为所有请求添加 Cors 中间件
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartsRouter(ApiGroup)

	return Router
}
