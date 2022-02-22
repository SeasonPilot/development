package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 使用 CORS 来允许跨源访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 服务端返回的 Access-Control-Allow-Origin: * 表明，该资源可以被 任意 外域访问。
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token") // 表明服务器允许请求中携带字段
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")                               // 表明服务器允许客户端使用哪些方法发起请求
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			// 并不是所有浏览器都支持预检请求的重定向。如果一个预检请求发生了重定向，一部分浏览器将报告错误：
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
