package handler

import (
	"context"
	"encoding/json"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsServer) GetAllCategoriesList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categories []model.Category
	// 预加载，查询出 二级 三级 分类
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categories)

	b, _ := json.Marshal(&categories)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}
