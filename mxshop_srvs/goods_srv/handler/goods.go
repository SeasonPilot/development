package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"

	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func ModelToGoodsInfoResp(good model.Goods) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              good.ID,
		CategoryId:      good.CategoryID,
		Name:            good.Name,
		GoodsSn:         good.GoodsSn,
		ClickNum:        good.ClickNum,
		SoldNum:         good.SoldNum,
		FavNum:          good.FavNum,
		MarketPrice:     good.MarketPrice,
		ShopPrice:       good.ShopPrice,
		GoodsBrief:      good.GoodsBrief,
		ShipFree:        good.ShipFree,
		GoodsFrontImage: good.GoodsFrontImage,
		IsNew:           good.IsNew,
		IsHot:           good.IsHot,
		OnSale:          good.OnSale,
		DescImages:      good.DescImages,
		Images:          good.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   good.Category.ID,
			Name: good.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   good.Brands.ID,
			Name: good.Brands.Name,
			Logo: good.Brands.Logo,
		},
	}
}

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func (s *GoodsServer) GoodsList(c context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	bq := elastic.NewBoolQuery()
	if req.PriceMin > 0 {
		bq.Filter(elastic.NewRangeQuery("shop_price").Gte(float32(req.PriceMin)))
	}
	if req.PriceMax > 0 {
		bq.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.IsHot {
		bq.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		bq.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.IsTab {
		bq.Filter(elastic.NewTermQuery("is_tab", req.IsTab))
	}
	// 顶级分类
	if req.TopCategory > 0 {
		// 验证顶级分类是否存在
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		var q string
		// 判断 顶级分类 Level
		if category.Level == 1 {
			q = fmt.Sprintf("select id from category where parent_category_id in (select id from category where parent_category_id = %d) ", category.ID)
		} else if category.Level == 2 {
			q = fmt.Sprintf("select id from category where parent_category_id = %d ", category.ID)
		} else if category.Level == 3 {
			q = fmt.Sprintf("select id from category where id = %d", category.ID)
		}

		type Result struct {
			ID int32
		}
		var results []Result
		// 在 MySQL 中查询所有 Category.ID
		global.DB.Model(&model.Category{}).Raw(q).Scan(&results)

		var categoryIDs []interface{}
		for _, r := range results {
			categoryIDs = append(categoryIDs, r.ID)
		}

		// 拿着查出来的所有 Category.ID 在 ES 中查询 所有 goodsIDs
		bq.Filter(elastic.NewTermsQuery("category_id", categoryIDs...))
	}

	if req.KeyWords != "" {
		bq.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.Brand > 0 {
		bq.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	// 执行 ES 搜索, 得到 goods
	goodsInES, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(bq).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(c)
	if err != nil {
		return nil, err
	}
	if len(goodsInES.Hits.Hits) == 0 {
		return nil, status.Errorf(codes.NotFound, "没有查询到相关商品")
	}

	rsp := new(proto.GoodsListResponse)

	// 要在分页前计算 total  fixme: 要再写一个不带分页 ES 查询语句吗？  这里 Total 应该一直是 req.PagePerNums ？？
	rsp.Total = int32(goodsInES.TotalHits())

	type Goods struct {
		ID int32 `json:"id"`
	}
	var (
		goodsES  Goods
		goodsIDs []int32
	)
	for _, hit := range goodsInES.Hits.Hits {
		err = json.Unmarshal(hit.Source, &goodsES)
		if err != nil {
			return nil, err
		}
		goodsIDs = append(goodsIDs, goodsES.ID)
	}

	var goods []model.Goods
	//fixme: Brands 预加载不出来,预加载全部 也不行
	// goods 结果中 BrandsID=0 ？？  为什么？？  model.Goods.BrandsID  映射的字段名称与数据库不一致，所以查出来的 BrandsID=0
	result := global.DB.Preload("Category").Preload("Brands").Find(&goods, goodsIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		rsp.Data = append(rsp.Data, ModelToGoodsInfoResp(good))
	}

	return rsp, nil
}

func (s *GoodsServer) BatchGetGoods(c context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Goods
	result := global.DB.Find(&goods, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.GoodsListResponse{}
	rsp.Total = int32(result.RowsAffected)

	for _, good := range goods {
		rsp.Data = append(rsp.Data, ModelToGoodsInfoResp(good))
	}

	return rsp, nil
}

func (s *GoodsServer) CreateGoods(c context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	goods := model.Goods{
		CategoryID:      req.CategoryId,
		Category:        category,
		BrandsID:        req.BrandId,
		Brands:          brand,
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
	}

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(c context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(c context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	goods = model.Goods{
		CategoryID:      req.CategoryId,
		Category:        category,
		BrandsID:        req.BrandId,
		Brands:          brand,
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
	}

	result := global.DB.Save(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) GetGoodsDetail(c context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	goods := model.Goods{}
	// fixme: 要返回关联表( category\brand )信息,要记得 Preload , Preload 中的字段是 model 中的
	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	return ModelToGoodsInfoResp(goods), nil
}
