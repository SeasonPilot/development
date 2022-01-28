package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

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

func removeTopStruct(fields map[string]string) map[string]string {
	resp := make(map[string]string, len(fields))
	for k, v := range fields {
		resp[k[strings.Index(k, ".")+1:]] = v
	}
	return resp
}

func initTrans(local string) error {
	v, ok := binding.Validator.Engine().(*validator.Validate) // 返回为默认 Validator 实例提供支持的底层验证器引擎。
	if !ok {
		return fmt.Errorf("validate err %v", v)
	}

	// 从tag中取值
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	enT := en.New() //中文翻译器
	zhT := zh.New() //英文翻译器
	//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
	universalTranslator := ut.New(enT, zhT, enT) // 通用翻译器

	var found bool
	// 根据 locale(区域) 返回指定的翻译器
	if trans, found = universalTranslator.GetTranslator(local); found {
		// 将 Translator 注册到 validator
		switch local {
		case "zh":
			_ = zh_translations.RegisterDefaultTranslations(v, trans) // 在 validator 中为所有内置标签注册一组默认翻译(Translations)
		case "en":
			_ = en_translations.RegisterDefaultTranslations(v, trans)
		default:
			_ = en_translations.RegisterDefaultTranslations(v, trans)
		}
		return nil
	}
	return fmt.Errorf("GetTranslator not found")
}

func main() {
	if err := initTrans("zh"); err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	r.POST("/loginJSON", func(c *gin.Context) {
		var loginForm LoginForm
		err := c.ShouldBind(&loginForm)
		if err != nil {
			// fixme:接口断言返回的是两个值
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusOK, gin.H{"err": errs.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"err": removeTopStruct(errs.Translate(trans))})
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
