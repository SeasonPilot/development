package router

import (
	"mxshop-api/order-web/api/shopcart"
	"mxshop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartsRouter(r *gin.RouterGroup) {
	group := r.Group("shopcart").Use(middlewares.JWTAuth())
	{
		group.GET("", shopcart.List)
		group.POST("", shopcart.New)
		//group.GET("/:id", shopcart.Details)
		group.DELETE("/:id", shopcart.Delete)
		group.PATCH("/:id", shopcart.Update)
		//group.PATCH("/:id", shopcart.Status)
	}
}
