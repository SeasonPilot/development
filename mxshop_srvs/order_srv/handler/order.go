package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	//更新购物车记录，更新数量和选中状态
	var cart model.ShoppingCart

	if result := global.DB.First(&cart, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	cart.Checked = request.Checked
	if request.Nums > 0 {
		cart.Nums = request.Nums
	}
	if result := global.DB.Save(&cart); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (OrderServer) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 删除时可以不用先查询，直接删除即可
	if result := global.DB.Delete(&model.ShoppingCart{}, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (OrderServer) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		1．从购物车中获取到选中的商品
		2．商品的价格自己查询—访问商品服务（跨微服务）
		3．库存的扣减—访问库存服务（跨微服务）
		4．订单的基本信息表—订单的商品信息表
		5．从购物车中删除已购买的记录
	*/

	// 1．从购物车中获取到选中的商品
	var (
		carts    []model.ShoppingCart
		goodsIDs []int32
	)
	// map 必现要初始化后才能使用
	goodsAndNums := make(map[int32]int32)

	if result := global.DB.Where(&model.ShoppingCart{User: request.UserId, Checked: true}).Find(&carts); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品")
	}

	for _, cart := range carts {
		goodsIDs = append(goodsIDs, cart.Goods)
		goodsAndNums[cart.Goods] = cart.Nums
	}

	// 2．商品的价格自己查询—访问商品服务（跨微服务）
	goods, err := global.GoodsClient.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{Id: goodsIDs})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var (
		totalPrice float32
		orderGoods []*model.OrderGoods // 这里是指针类型, 是后面哪里要修改这个变量？？？
		goodsInfo  []*proto.GoodsInvInfo
	)
	for _, good := range goods.Data {
		totalPrice += good.ShopPrice * float32(goodsAndNums[good.Id])

		orderGoods = append(orderGoods, &model.OrderGoods{
			Order:      request.Id, // 这里不传？？？
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsAndNums[good.Id], // 怎么从 carts 中拿
		})

		goodsInfo = append(goodsInfo, &proto.GoodsInvInfo{
			GoodsID: good.Id,
			Num:     goodsAndNums[good.Id],
		})
	}

	// 4．订单的基本信息表—订单的商品信息表
	orderInfo := model.OrderInfo{
		User:         request.UserId,
		OrderMount:   totalPrice,
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
	}

	global.DB.Save(&orderGoods)
	global.DB.Save(&orderInfo)

	// 3．库存的扣减—访问库存服务（跨微服务）
	_, err = global.InventoryClient.Sell(ctx, &proto.SellInfo{
		GoodsInfo: goodsInfo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// 5．从购物车中删除已购买的记录     可不可以调用微服务自己的方法 DeleteCartItem ？？？
	if result := global.DB.Delete(&carts); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &proto.OrderInfoResponse{
		UserId:  0,
		OrderSn: "",
		PayType: "",
		Status:  "",
		Total:   totalPrice,
		AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (OrderServer) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var (
		total  int64
		rsp    proto.OrderListResponse
		orders []model.OrderInfo
	)

	global.DB.Where(model.OrderInfo{User: request.UserId}).Count(&total)
	rsp.Total = int32(total)

	global.DB.Scopes(Paginate(int(request.Pages), int(request.PagePerNums))).Where(&model.OrderInfo{User: request.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (OrderServer) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var (
		order model.OrderInfo
		rsp   proto.OrderInfoDetailResponse
		goods []model.OrderGoods
	)

	//这个订单的id是否是当前用户的订单， 如果在web层用户传递过来一个id的订单， web层应该先查询一下订单id是否是当前用户的
	//在个人中心可以这样做，但是如果是后台管理系统，web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	if result := global.DB.Where(&model.OrderInfo{User: request.UserId, BaseModel: model.BaseModel{ID: request.Id}}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	if result := global.DB.Find(&goods, "order", request.Id); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range goods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (OrderServer) UpdateOrderStatus(ctx context.Context, status *proto.OrderStatus) (*emptypb.Empty, error) {
	panic("implement me")
}
