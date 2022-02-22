package router

import (
	"mxshop-api/goods-web/api/goods"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	group := r.Group("good")
	{
		group.GET("list", goods.List)
	}
}
