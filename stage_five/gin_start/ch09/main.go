package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("templates/goods.tmpl")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.tmpl", gin.H{
			"title": "hello season",
		})
	})

	r.Run()
}
