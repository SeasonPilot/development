package pay

import (
	"net/http"

	"mxshop-api/order-web/api"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

// Notify 支付宝回调通知
func Notify(c *gin.Context) {
	client, err := alipay.New(global.SrvConfig.AlipayInfo.AppID, global.SrvConfig.AlipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.SrvConfig.AlipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 验证签名
	noti, err := client.GetTradeNotification(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// 更新订单状态
	_, err = global.OrderClient.UpdateOrderStatus(c, &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		zap.S().Error("更新订单状态失败")
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.String(http.StatusOK, "success")
}
