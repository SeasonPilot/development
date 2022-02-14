package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 入参是函数调用
	r.Use(MyLogger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}

// MyLogger 没有入参，返回值是 HandlerFunc
func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Set("fa", "fsf")
		//让原本该执行的逻辑继续执行
		c.Next()

		fmt.Printf("耗时:%v\n", time.Since(now))
		fmt.Println("状态", c.Writer.Status())
	}
}
