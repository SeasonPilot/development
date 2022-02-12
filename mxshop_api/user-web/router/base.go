package router

import (
	"development/mxshop_api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {
	group := r.Group("base")
	{
		group.GET("captcha", api.GetCaptcha)
	}
}
