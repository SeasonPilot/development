package router

import (
	"development/mxshop_api/user-web/api"
	"development/mxshop_api/user-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	// user 前是否一定要 /   不是
	group := r.Group("user")
	{
		group.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		group.POST("login", api.PassWordLogin)
		group.POST("register", api.Register)
	}
}
