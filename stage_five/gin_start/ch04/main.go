package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	ID     int    `uri:"id" binding:"required"`
	Action string `uri:"action" binding:"required"`
}

func main() {
	r := gin.Default()
	group := r.Group("/goods")
	{
		group.GET("/:id/:action", GetGoods)
		//group.GET("/:id/*action", GetGoods) // 匹配 action 后所有的路径，常用于文件地址
		group.POST("/create")
	}
	r.Run()
}

func GetGoods(c *gin.Context) {
	//id := c.Param("id")
	//action := c.Param("action")
	var person Person
	err := c.ShouldBindUri(&person)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"id":     person.ID,
		"action": person.Action,
	})
}
