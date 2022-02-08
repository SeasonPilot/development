package router

import (
	"development/mxshop_api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	// user 前是否一定要 /   不是
	group := r.Group("user")
	{
		group.GET("list", api.GetUserList)
	}
}
