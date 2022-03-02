package handler

import (
	"context"

	"mxshop-srvs/order_srv/global"
	"mxshop-srvs/order_srv/model"
	"mxshop-srvs/order_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func (OrderServer) CartItemList(ctx context.Context, info *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var (
		carts []model.ShoppingCart
		rsp   proto.CartItemListResponse
	)

	if result := global.DB.Where(&model.ShoppingCart{User: info.Id}).Find(&carts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, cart := range carts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			// fixme: 不是 false
			Checked: cart.Checked,
		})
	}

	return &rsp, nil
}

func (OrderServer) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车- 合并
	var cart model.ShoppingCart

	// fixme: 购物车商品记录是属于用户的，所以查询的时候要加上用户ID
	if result := global.DB.First(&cart, "goods = ? AND user = ?", request.GoodsId, request.UserId); result.RowsAffected == 1 {
		//如果记录已经存在，则合并购物车记录, 更新操作
		cart.Nums += request.Nums
	} else {
		//插入操作
		cart.Nums = request.Nums
		cart.User = request.UserId
		cart.Goods = request.GoodsId
		cart.Checked = false
	}

	if result := global.DB.Save(&cart); result.Error != nil {
		return nil, result.Error
	}

	// 只返回 ID 即可,其他都是前端传入的,前端是知道的
	return &proto.ShopCartInfoResponse{Id: cart.ID}, nil
}

func (OrderServer) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (OrderServer) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (OrderServer) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	panic("implement me")
}

func (OrderServer) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	panic("implement me")
}

func (OrderServer) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	panic("implement me")
}

func (OrderServer) UpdateOrderStatus(ctx context.Context, status *proto.OrderStatus) (*emptypb.Empty, error) {
	panic("implement me")
}
