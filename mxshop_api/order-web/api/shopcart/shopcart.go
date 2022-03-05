package shopcart

import (
	"fmt"
	"net/http"

	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	// 从用户上下文读取值
	// 从 Context 中获取当前登陆的用户
	userId, ok := c.Get("userId")
	if !ok {
		fmt.Println(userId, ok)
		c.Status(http.StatusBadRequest)
		return
	}

	cartListRsp, err := global.OrderClient.CartItemList(c, &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败: %s", err.Error())
		api.RpcErrToHttpErr(c, err)
		return
	}

	var goodsIDs []int32
	for _, item := range cartListRsp.Data {
		goodsIDs = append(goodsIDs, item.GoodsId)
	}
	if len(goodsIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	goodsListResp, err := global.GoodsClient.BatchGetGoods(c, &proto.BatchGoodsIdInfo{Id: goodsIDs})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	goodsRsp := make([]interface{}, 0)
	for _, item := range cartListRsp.Data {
		for _, good := range goodsListResp.Data {
			if good.Id == item.GoodsId {
				tmpMap := make(map[string]interface{})
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = good.Name
				tmpMap["good_image"] = good.GoodsFrontImage
				tmpMap["good_price"] = good.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsRsp = append(goodsRsp, tmpMap)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": cartListRsp.Total,
		"data":  goodsRsp,
	})
}

func New(c *gin.Context) {
	var itemForm forms.ShopCartItemForm
	err := c.ShouldBind(&itemForm)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	// 检查商品是否存在
	_, err = global.GoodsClient.GetGoodsDetail(c, &proto.GoodInfoRequest{Id: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[List] 查询【商品信息】失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	// 校验库存是否足够
	invDetail, err := global.InvClient.InvDetail(c, &proto.GoodsInvInfo{GoodsID: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[List] 查询【库存信息】失败")
		api.RpcErrToHttpErr(c, err)
		return
	}
	if invDetail.Num < itemForm.Num {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("库存不足，仅剩 %d 件", invDetail.Num),
		})
		return
	}

	// userId 从 context 中拿，不是让前端传入
	userId, _ := c.Get("userId")
	item, err := global.OrderClient.CreateCartItem(c, &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: itemForm.GoodsId,
		Nums:    itemForm.Num,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": item.Id,
	})
}

func Details(c *gin.Context) {

}
func Delete(c *gin.Context) {

}
func Update(c *gin.Context) {

}
func Status(c *gin.Context) {

}
