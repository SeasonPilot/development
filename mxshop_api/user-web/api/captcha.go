package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.DefaultDriverDigit

	// 创建 Captcha 对象
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误：%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captcha_id": id,
		"captcha":    b64s,
	})
}
