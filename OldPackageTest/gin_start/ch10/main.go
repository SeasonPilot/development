package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

/*
	优雅退出，当我们关闭程序的时候应该做的后续处理
	微服务启动之前或者启动之后会做一件事: 将当前的服务的ip地址和端号注册到注册中心
	我们当前的服务停止了以后并没有告知注册中心
*/
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go func() {
		r.Run()
	}()

	// 在 r.Run() 退出前做一些逻辑
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	fmt.Println("关闭 server 中。。。")
	fmt.Println("注销服务。。。")
}
