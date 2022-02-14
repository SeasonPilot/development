package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	group := r.Group("/goods")
	{
		group.GET("list", listGoods)
		group.POST("create", CreateGoods)
	}
}

func CreateGoods(context *gin.Context) {}

func listGoods(context *gin.Context) {}
