package api

import (
	"net/http"
	"strings"

	"mxshop-api/goods-web/global"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
removeTopStruct 移除 `PassWordLoginForm.`
{
    "msg": {
        "PassWordLoginForm.password": "password长度必须至少为3个字符"
    }
}
*/
func removeTopStruct(fields map[string]string) map[string]string {
	resp := make(map[string]string, len(fields))
	for k, v := range fields {
		resp[k[strings.Index(k, ".")+1:]] = v
	}
	return resp
}

// RpcErrToHttpErr 将 grpc 的 code 转换成 http 的状态码
func RpcErrToHttpErr(c *gin.Context, err error) {
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			switch grpcStatus.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": grpcStatus.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					// 不能直接返回  grpcStatus.Message() ,会暴露过多信息给用户,如敏感信息。
					// 不能把 grpc 内部错误暴露给用户，不友好
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "商品服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": grpcStatus.Code(),
				})
			}
		}
	}
}

// HandleValidatorError 处理 Validator 的错误
func HandleValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": removeTopStruct(errs.Translate(global.Translator)),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": err.Error(),
	})
}
