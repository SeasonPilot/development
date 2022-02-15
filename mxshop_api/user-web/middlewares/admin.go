package middlewares

import (
	"net/http"

	"mxshop-api/user-web/models"

	"github.com/gin-gonic/gin"
)

// IsAdminAuth 验证用户是否为管理员
func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		currentUser, ok := claims.(*models.CustomClaims)
		if ok {
			if currentUser.AuthorityID == 2 {
				c.Next()
				// fixme: c.Next() 后面也要 return
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "没有权限",
		})
		c.Abort()
	}
}
