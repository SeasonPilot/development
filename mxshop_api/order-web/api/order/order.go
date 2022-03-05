package order

import (
	"net/http"
	"strconv"

	"mxshop-api/order-web/api"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/models"
	"mxshop-api/order-web/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*models.CustomClaims)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	req := &proto.OrderFilterRequest{}
	if customClaims.AuthorityID == 1 {
		req.UserId = int32(userId.(uint))
	}

	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Pages = int32(pagesInt)

	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNums = int32(perNumsInt)

	orderListResp, err := global.OrderClient.OrderList(c, req)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	rspMap := make([]interface{}, 0)
	for _, item := range orderListResp.Data {
		tmpMap := make(map[string]interface{})

		tmpMap["userId"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["add_time"] = item.AddTime

		rspMap = append(rspMap, tmpMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"total": orderListResp.Total,
		"data":  rspMap,
	})
}

func New(c *gin.Context) {

}

func Details(c *gin.Context) {

}
func Delete(c *gin.Context) {

}
func Update(c *gin.Context) {

}
