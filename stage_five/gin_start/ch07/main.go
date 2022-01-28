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

type SignUpForm struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"` // uint8 类型
	Name       string `json:"name"  binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
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

	r.POST("/signUp", func(c *gin.Context) {
		err := c.ShouldBind(&SignUpForm{})
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	_ = r.Run()
}
