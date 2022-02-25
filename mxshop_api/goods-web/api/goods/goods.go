package goods

import (
	"net/http"
	"strconv"

	"mxshop-api/goods-web/api"
	"mxshop-api/goods-web/forms"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"mxshop-api/goods-web/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  response.RespToModels(rsp.Data),
	})
}

func New(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	err := c.ShouldBind(&goodsForm)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	goods, err := global.GoodsClient.CreateGoods(c, &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, response.RespToModels([]*proto.GoodsInfoResponse{goods}))
}

func Details(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	rsp, err := global.GoodsClient.GetGoodsDetail(c, &proto.GoodInfoRequest{
		Id: int32(idInt),
	})
	if err != nil {
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, response.RespToModels([]*proto.GoodsInfoResponse{rsp}))
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	_, err = global.GoodsClient.DeleteGoods(c, &proto.DeleteGoodsInfo{
		Id: int32(idInt),
	})
	if err != nil {
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func Stocks(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	//TODO 商品的库存
	c.JSON(http.StatusOK, nil)
}

func UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	goodsStatusForm := forms.GoodsStatusForm{}
	err = c.ShouldBind(&goodsStatusForm)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	_, err = global.GoodsClient.UpdateGoods(c, &proto.CreateGoodsInfo{
		Id:     int32(idInt),
		IsNew:  *goodsStatusForm.IsNew,
		IsHot:  *goodsStatusForm.IsHot,
		OnSale: *goodsStatusForm.OnSale,
	})
	if err != nil {
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	goodsForm := forms.GoodsForm{}
	err = c.ShouldBind(&goodsForm)
	if err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	_, err = global.GoodsClient.UpdateGoods(c, &proto.CreateGoodsInfo{
		Id:              int32(idInt),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.RpcErrToHttpErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
