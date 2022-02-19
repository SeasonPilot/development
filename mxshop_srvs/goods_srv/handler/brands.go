package handler

import (
	"context"

	"mxshop-srvs/goods_srv/global"
	"mxshop-srvs/goods_srv/model"
	"mxshop-srvs/goods_srv/proto"
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
