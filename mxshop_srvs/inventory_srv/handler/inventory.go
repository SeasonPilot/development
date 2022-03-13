package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"mxshop-srvs/inventory_srv/global"
	"mxshop-srvs/inventory_srv/model"
	"mxshop-srvs/inventory_srv/proto"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
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
//var m sync.Mutex

func (InventoryServer) Sell(ctx context.Context, info *proto.SellInfo) (*emptypb.Empty, error) {
	// 一个订单是一个事务,一个订单(购物车)包含多个商品
	tx := global.DB.Begin()
	//m.Lock()

	var goodsDetailList model.GoodsDetailList

	for _, good := range info.GoodsInfo {
		mutex := global.RedSync.NewMutex(fmt.Sprintf("goods_%d", good.GoodsID))

		err := mutex.Lock()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常: %s", err.Error())
		}

		var inv model.Inventory
		if result := global.DB.First(&inv, "Goods = ?", good.GoodsID); result.RowsAffected == 0 {
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
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}

		ok, err := mutex.Unlock()
		if !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常: %s", err.Error())
		}

		goodsDetailList = append(goodsDetailList, model.GoodsDetail{
			Goods: good.GoodsID,
			Num:   good.Num,
		})
	}

	// 保存库存销售记录
	var stockSellDetail = model.StockSellDetail{
		OrderSN: info.OrderSn,
		Status:  1,
		Detail:  goodsDetailList,
	}
	if result := tx.Create(&stockSellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存销售记录失败: %s", result.Error.Error())
	}

	tx.Commit()
	//m.Unlock()
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

func AutoReback(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSN string
	}
	var orderInfo OrderInfo

	tx := global.DB.Begin()

	for _, msg := range ext {
		if err := json.Unmarshal(msg.Body, &orderInfo); err != nil {
			zap.S().Errorf("解析json失败： %v\n", msg.Body)
			return consumer.ConsumeSuccess, nil
		}

		var stockSellDetail model.StockSellDetail

		// 查询扣减记录是否存在
		if result := tx.Where(model.StockSellDetail{
			OrderSN: orderInfo.OrderSN,
			Status:  1,
		}).First(&stockSellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}

		// 如果查询到那么逐个归还库存
		for _, detail := range stockSellDetail.Detail {
			var inv model.Inventory

			if result := tx.Where(&inv).Update("stocks", gorm.Expr("stocks + ?", detail.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
	}

	// 更新库存销售记录扣减状态
	if result := tx.Where(model.StockSellDetail{OrderSN: orderInfo.OrderSN}).Update("status", 2); result.RowsAffected == 0 {
		tx.Rollback()
		return consumer.ConsumeRetryLater, nil
	}

	tx.Commit()
	return consumer.ConsumeSuccess, nil
}
