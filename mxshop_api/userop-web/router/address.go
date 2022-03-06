package router

import (
	"mxshop-api/userop-web/api/address"
	"mxshop-api/userop-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	{
		AddressRouter.GET("", address.List)          // 列表页
		AddressRouter.DELETE("/:id", address.Delete) // 删除
		AddressRouter.POST("", address.New)          //新建
		AddressRouter.PUT("/:id", address.Update)    //修改
	}
}
