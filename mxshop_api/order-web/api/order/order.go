package order

import (
	"net/http"
	"strconv"

	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
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
	orderForm := forms.CreateOrderForm{}
	err := c.ShouldBind(&orderForm)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	userId, _ := c.Get("userId")
	rsp, err := global.OrderClient.CreateOrder(c, &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Address: orderForm.Addr,
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Details(c *gin.Context) {
	// c.Param("id") 入参要与 router 中 /:id 一致
	orderId := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "订单 ID 参数错误",
		})
		return
	}

	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")

	customClaims, ok := claims.(*models.CustomClaims)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	req := &proto.OrderRequest{}
	if customClaims.AuthorityID == 1 {
		req.UserId = int32(userId.(uint))
	}
	req.Id = int32(orderIdInt)

	rsp, err := global.OrderClient.OrderDetail(c, req)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	goods := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := map[string]interface{}{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goods = append(goods, tmpMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       rsp.OrderInfo.Id,
		"status":   rsp.OrderInfo.Status,
		"user":     rsp.OrderInfo.UserId,
		"post":     rsp.OrderInfo.Post,
		"total":    rsp.OrderInfo.Total,
		"address":  rsp.OrderInfo.Address,
		"name":     rsp.OrderInfo.Name,
		"mobile":   rsp.OrderInfo.Mobile,
		"pay_type": rsp.OrderInfo.PayType,
		"order_sn": rsp.OrderInfo.OrderSn,
		"goods":    goods,
	})
}
