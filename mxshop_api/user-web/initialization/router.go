package initialization

import (
	"development/mxshop_api/user-web/router"

	"github.com/gin-gonic/gin"
)

// Routers 全局Routers; 避免每个 router(如 user、短信等) 都实例化一个
func Routers() *gin.Engine {
	// 大写是为了和包名区分开
	Router := gin.Default()

	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}
