package router

import (
	"mxshop-api/goods-web/api/goods"
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	group := r.Group("good")
	{
		group.GET("", middlewares.JWTAuth(), goods.List)
		group.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		group.GET("/:id", goods.Details)
		group.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		group.GET("/:id/stocks", goods.Stocks)

		group.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		group.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
}
