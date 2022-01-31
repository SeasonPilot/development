package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 加载二级目录下所有文件
	r.LoadHTMLGlob("templates/**/*")
	//r.LoadHTMLFiles("templates/goods.tmpl")

	// 第一个参数是 url 。第二个参数是静态文件目录,即当前目录下的 static 文件夹
	r.Static("/static", "./static")

	r.GET("/index", func(c *gin.Context) {
		//如果没有在模板中使用define定义 那么我们就可以使用默认的文件名来找
		c.HTML(http.StatusOK, "default/index", gin.H{
			"title": "hello season",
		})
	})

	r.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods/list", gin.H{
			"title": "商品列表",
		})
	})

	r.GET("users/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/list", gin.H{
			"title": "商品列表",
		})
	})

	r.Run()
}
