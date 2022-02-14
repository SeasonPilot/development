package main

import (
	"development/OldPackageTest/gin_start/ch06/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/moreJSON", moreJson)
	r.GET("/proto", someProto)

	r.Run()
}

// 返回 protobuf
func someProto(c *gin.Context) {
	teach := proto.Teacher{
		Name:   "bobby",
		Course: []string{"go", "python", "微服务"},
	}
	c.ProtoBuf(http.StatusOK, &teach)
}

// josn 直接返回结构体
func moreJson(c *gin.Context) {
	type Msg struct {
		Name    string `json:"user"`
		Message string
		Number  int
	}

	msg := Msg{
		Name:    "hhh",
		Message: "这是一个测试json",
		Number:  121,
	}
	c.JSON(http.StatusOK, &msg)
}
