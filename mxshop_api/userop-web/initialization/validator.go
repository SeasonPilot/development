package initialization

import (
	"reflect"
	"strings"

	"mxshop-api/userop-web/global"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(local string) {
	//1. validator engine 断言成 Validate
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		//2. 注册一个函数来获取 StructFields 的备用名称。并放在 Validate 中
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 不要自己 new Validate, 要用 binding.Validator.Engine() 断言的 Validate
		// validate := validator.New()

		//3. 通用翻译器
		enT := en.New() //中文翻译器
		zhT := zh.New() //英文翻译器
		// 第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		universalTranslator := ut.New(enT, zhT, enT) // 通用翻译器

		//4. 根据形参 GetTranslator 拿到具体的翻译器
		var found bool
		global.Translator, found = universalTranslator.GetTranslator(local)
		if found {
			//5.  将 Translator 翻译器 注册到 Validate
			switch local {
			case "en":
				err := en2.RegisterDefaultTranslations(v, global.Translator)
				if err != nil {
					return
				}
			case "zh":
				err := zh2.RegisterDefaultTranslations(v, global.Translator)
				if err != nil {
					return
				}
			default:
				err := en2.RegisterDefaultTranslations(v, global.Translator)
				if err != nil {
					return
				}
			}
		}
	}
}
