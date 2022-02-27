package handler

import (
	"context"

	"mxshop-srvs/inventory_srv/global"
	"mxshop-srvs/inventory_srv/model"
	"mxshop-srvs/inventory_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// SetInv 既可以新增也可以更新
func (InventoryServer) SetInv(ctx context.Context, info *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.First(&inv, "Goods = ?", info.GoodsID)
	inv.Goods = info.GoodsID
	inv.Stocks = info.Num

	if result := global.DB.Save(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (InventoryServer) InvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.First(&inv, "Goods = ?", info.GoodsID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品库存信息不存在")
	}

	return &proto.GoodsInvInfo{
		GoodsID: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (InventoryServer) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	panic("implement me")
}

func (InventoryServer) Reback(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	panic("implement me")
}
