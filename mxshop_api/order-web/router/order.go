package router

import (
	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/api/pay"
	"mxshop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(r *gin.RouterGroup) {
	group := r.Group("order").Use(middlewares.JWTAuth())
	{
		group.GET("", order.List)
		group.POST("", order.New)
		group.GET("/:id", order.Details)
	}
	payGroup := r.Group("pay")
	{
		payGroup.POST("alipay/notify", pay.Notify)
	}
}
