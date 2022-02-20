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
	// 嵌套预加载  补充 SQL 语句
	// SELECT * FROM `category` WHERE `category`.`level` = 1
	// SELECT * FROM `category` WHERE `category`.`parent_category_id` IN (130358,130361,135200,135201,135202)
	// SELECT * FROM `category` WHERE `category`.`parent_category_id` IN (130364,130365,135486,135487,135488,135489,130370,136604,136614,136624,136634,136643,136654,136661,136669,136678,136688,136698,136708,136719,136728,136733,136742,136751,146239,136760,136770,136781,136789,136800,136805,150564,164593,169166,136818,136828,136838,136848,136854,136865,136870,136880,136891)

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

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}

	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab
	if req.Level != 1 {
		//去查询父类目是否存在
		category.ParentCategoryID = req.ParentCategory
	}
	global.DB.Save(&category)

	return &proto.CategoryInfoResponse{Id: category.ID}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	// 有问题: bool 零值为 false; 如果用户想传 false 就传不进来
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
