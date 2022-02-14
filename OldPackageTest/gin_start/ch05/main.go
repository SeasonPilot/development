package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "2")
		name := c.PostForm("name")
		age := c.DefaultPostForm("age", "20")
		c.JSON(http.StatusOK, gin.H{
			"id":   id,
			"page": page,
			"name": name,
			"age":  age,
		})
	})

	router.Run()
}
