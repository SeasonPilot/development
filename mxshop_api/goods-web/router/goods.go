package router

import (
	"mxshop-api/goods-web/api/goods"
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	group := r.Group("good")
	{
		group.GET("", goods.List)
		group.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		group.GET(":id", goods.Details)
	}
}
