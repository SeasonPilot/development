package initialization

import (
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

// Routers 全局Routers; 避免每个 router(如 user、base等) 都实例化一个
func Routers() *gin.Engine {
	// 大写是为了和包名区分开
	Router := gin.Default()

	// 为所有请求添加 Cors 中间件
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
