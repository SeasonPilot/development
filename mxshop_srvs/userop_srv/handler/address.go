package handler

import (
	"context"

	"mxshop-srvs/userop_srv/global"
	"mxshop-srvs/userop_srv/model"
	"mxshop-srvs/userop_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	var rsp proto.AddressListResponse
	var addressResponse []*proto.AddressResponse

	if result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addresses); result.RowsAffected != 0 {
		rsp.Total = int32(result.RowsAffected)
	} else if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取地址错误: %s", result.Error)
	}

	for _, address := range addresses {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	rsp.Data = addressResponse

	return &rsp, nil
}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	var address model.Address

	address.User = req.UserId
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.Address = req.Address
	address.SignerName = req.SignerName
	address.SignerMobile = req.SignerMobile

	if result := global.DB.Save(&address); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址错误: %s", result.Error)
	}

	return &proto.AddressResponse{Id: address.ID}, nil
}

func (*UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).Delete(&model.Address{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	var address model.Address

	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if address.Province != "" {
		address.Province = req.Province
	}

	if address.City != "" {
		address.City = req.City
	}

	if address.District != "" {
		address.District = req.District
	}

	if address.Address != "" {
		address.Address = req.Address
	}

	if address.SignerName != "" {
		address.SignerName = req.SignerName
	}

	if address.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}

	global.DB.Save(&address)

	return &emptypb.Empty{}, nil
}
