package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	User     string `form:"user" json:"user" binding:"required,min=3,max=10"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	r.POST("/loginJSON", func(c *gin.Context) {
		var loginForm LoginForm
		err := c.ShouldBind(&loginForm)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "登陆成功"})
	})

	_ = r.Run()
}
