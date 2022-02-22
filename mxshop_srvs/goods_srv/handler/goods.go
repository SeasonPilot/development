package handler

import (
	"context"
	"fmt"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"

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
	localDB := global.DB.Model(&model.Goods{})

	if req.PriceMin > 0 {
		// 不重新赋值 localDB 好像也没有问题？？
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}
	if req.IsHot {
		localDB = localDB.Where("is_hot = true")
	}
	if req.IsNew {
		localDB = localDB.Where("is_new = true")
	}
	if req.IsTab {
		localDB = localDB.Where("is_tab = true")
	}
	// 顶级分类
	if req.TopCategory > 0 {
		// 验证顶级分类是否存在
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		// 判断 顶级分类 Level
		if category.Level == 1 {
			// 查询 1 级分类 包含的所有商品
			localDB = localDB.Where(fmt.Sprintf("category_id in ( select id from category where parent_category_id in (select id from category where parent_category_id = %d) )", category.ID))
		} else if category.Level == 2 {
			// 查询 2 级分类 包含的所有商品
			localDB = localDB.Where(fmt.Sprintf("category_id in (select id from category where parent_category_id in (%d) )", category.ID))
		} else if category.Level == 3 {
			// 查询 3 级分类 包含的所有商品
			localDB = localDB.Where(fmt.Sprintf("category_id in (%d)", category.ID))
		}
	}

	if req.KeyWords != "" {
		// 字段名称要与数据库中的 column 一致，不是 struct 字段,LIKE 后面要加 %
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}

	rsp := new(proto.GoodsListResponse)

	// 要在分页前计算 total
	var count int64
	localDB.Count(&count)
	rsp.Total = int32(count)

	var goods []model.Goods
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
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

	result := global.DB.Save(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

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
