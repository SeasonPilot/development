package router

import (
	"mxshop-api/userop-web/api/message"
	"mxshop-api/userop-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 列表页
		MessageRouter.POST("", message.New) //新建
	}
}
