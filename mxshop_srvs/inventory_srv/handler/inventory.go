package handler

import (
	"context"
	"sync"

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

// 全局锁,所有协程共用一把锁
var m sync.Mutex

func (InventoryServer) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 一个订单是一个事务,一个订单(购物车)包含多个商品
	tx := global.DB.Begin()
	m.Lock()

	for _, good := range info.GoodsInfo {
		var inv model.Inventory
		if result := tx.First(&inv, "Goods = ?", good.GoodsID); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "GoodsID: %d, 商品库存信息不存在", good.GoodsID)
		}

		// 判断库存是否充足
		if inv.Stocks < good.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "GoodsID: %d, 库存不足", good.GoodsID)
		}

		// 扣减库存
		inv.Stocks -= good.Num // 要记得扣除库存
		if result := tx.Save(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
	}

	tx.Commit()
	m.Unlock()
	return &emptypb.Empty{}, nil
}

func (InventoryServer) Reback(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()

	for _, good := range info.GoodsInfo {
		var inv model.Inventory
		if result := tx.First(&inv, "Goods = ?", good.GoodsID); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "GoodsID: %d, 商品库存信息不存在", good.GoodsID)
		}

		// 归还库存
		inv.Stocks += good.Num
		if result := tx.Save(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}
