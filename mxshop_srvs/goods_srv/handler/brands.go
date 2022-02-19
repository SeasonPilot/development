package handler

import (
	"context"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *GoodsServer) BrandList(c context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	resp := proto.BrandListResponse{}

	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var count int64
	global.DB.Model(&brands).Count(&count)
	resp.Total = int32(count)

	var data []*proto.BrandInfoResponse
	for _, brand := range brands {
		data = append(data, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	resp.Data = data
	return &resp, nil
}

func (s *GoodsServer) CreateBrand(c context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	if result := global.DB.First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Error(codes.InvalidArgument, "品牌已存在")
	}

	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	if result := global.DB.Save(&brand); result.Error != nil {
		return nil, status.Error(codes.Internal, "保存 brand 失败")
	}

	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}
func (s *GoodsServer) DeleteBrand(c context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "品牌不存在")
		}
		return nil, status.Error(codes.Internal, "删除 brand 失败")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBrand(c context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand := model.Brands{}
	if result := global.DB.First(&brand, req.Id); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "品牌不存在")
	}

	// 只有传入参数不为空是才更新字段，防止用户不小心传入空值
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}

	if result := global.DB.Save(&brand); result.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "保存 brand 失败")
	}

	return &emptypb.Empty{}, nil
}
