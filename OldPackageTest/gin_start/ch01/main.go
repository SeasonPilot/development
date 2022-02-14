package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//实例化一个gin的server对象
	g := gin.Default()
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := g.Run()
	if err != nil {
		panic(err)
	}
}
