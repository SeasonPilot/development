package goods

import (
	"net/http"
	"strconv"
	"strings"

	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"mxshop-api/goods-web/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
func RpcErrToHttpErr(err error, c *gin.Context) {
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
					"msg": "用户服务不可用",
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

func List(c *gin.Context) {
	var req proto.GoodsFilterRequest

	pMax := c.DefaultQuery("pmax", "0")
	pMaxInt, _ := strconv.Atoi(pMax)
	req.PriceMax = int32(pMaxInt)

	pMin := c.DefaultQuery("pmin", "0")
	pMinInt, _ := strconv.Atoi(pMin)
	req.PriceMin = int32(pMinInt)

	isHot := c.DefaultQuery("ih", "0")
	// 应该是字符串 1 ,且不用 Atoi
	if isHot == "1" {
		req.IsHot = true
	}

	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		req.IsNew = true
	}

	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		req.IsNew = true
	}

	category := c.DefaultQuery("category", "0")
	categoryInt, _ := strconv.Atoi(category)
	req.TopCategory = int32(categoryInt)

	kw := c.DefaultQuery("keyword", "")
	req.KeyWords = kw

	brand := c.DefaultQuery("brand", "0")
	brandInt, _ := strconv.Atoi(brand)
	req.Brand = int32(brandInt)

	// 忘记分页了
	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Pages = int32(pagesInt)

	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNums = int32(perNumsInt)

	rsp, err := global.GoodsClient.GoodsList(c, &req)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		RpcErrToHttpErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  response.RespToModels(rsp.Data),
	})
}
