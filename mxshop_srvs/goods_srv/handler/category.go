package handler

import (
	"context"
	"encoding/json"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GoodsServer) GetAllCategoriesList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categories []model.Category
	// 预加载，查询出 二级 三级 分类
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categories)

	b, _ := json.Marshal(&categories)

	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

func (s *GoodsServer) GetSubCategory(c context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 判断目标是否存在
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}

	// 指针类型的变量要初始化
	resp := new(proto.SubCategoryListResponse)
	// 返回目标分类信息
	resp.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	var subCategories []*model.Category
	var subCategoriesInfoResp []*proto.CategoryInfoResponse

	// 查找目标 下一级 level 的分类
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategories)
	// 查找目标 下一级 level 和 下两级 level 的分类; 这种用 JSON 格式返回比较容易，不然要拼凑 slice
	//global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload("SubCategory").Find(&subCategories)

	for _, subCategory := range subCategories {
		subCategoriesInfoResp = append(subCategoriesInfoResp, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})

	}

	resp.SubCategories = subCategoriesInfoResp
	return resp, nil
}
