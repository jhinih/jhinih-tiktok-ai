package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/repository"
	"Tiktok/response"

	"Tiktok/types"
	"Tiktok/utils"
	"context"
	"errors"
	"strconv"
	"time"
)

type ShopLogic struct {
}

func NewShopLogic() *ShopLogic {
	return &ShopLogic{}
}

func (l *ShopLogic) GetOrderInfo(ctx context.Context, req types.GetOrderInfoRequest) (resp types.GetOrderInfoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// 从数据库获取订单列表
	OrderID, _ := strconv.ParseInt(req.OrderID, 10, 64)
	Order, err := repository.NewShopRequest(global.DB).GetOrderInfo(OrderID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频列表失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	resp.Order = types.Order{
		ID:          strconv.FormatInt(Order.ID, 10),
		CreatedTime: strconv.FormatInt(Order.CreatedTime, 10),
		UpdatedTime: strconv.FormatInt(Order.UpdatedTime, 10),
		OwnerID:     strconv.FormatInt(Order.OwnerID, 10),
		IsDelete:    Order.IsDelete,
		Price:       strconv.FormatInt(Order.Price, 10),
		Status:      strconv.Itoa(Order.Status),
		Pay:         Order.Pay,
		PayBy:       Order.PayBy,
		PayTime:     strconv.FormatInt(Order.PayTime, 10),
	}
	return resp, nil
}
func (l *ShopLogic) CreateOrder(ctx context.Context, req types.CreateOrderRequest) (resp types.CreateOrderResponse, err error) {
	defer utils.RecordTime(time.Now())()
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	Price, _ := strconv.ParseInt(req.Price, 10, 64)
	OrderID := global.SnowflakeNode.Generate().Int64()
	Order := model.Order{
		ID:        OrderID,
		TimeModel: model.TimeModel{},
		OwnerID:   OwnerID,
		IsDelete:  false,
		Price:     Price,
		Status:    0,
		Pay:       false,
		PayBy:     "",
		PayTime:   0,
	}

	err = repository.NewShopRequest(global.DB).CreateOrder(Order)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存订单到数据库失败: %v", err)
		return types.CreateOrderResponse{}, errors.New("保存订单信息失败")
	}
	zlog.CtxInfof(ctx, "订单信息保存到数据库成功, ID: %d", OrderID)

	zlog.CtxInfof(ctx, "订单保存数据库处理完成: %+v", Order)
	return resp, nil
}
func (l *ShopLogic) UpdateOrder(ctx context.Context, req types.UpdateOrderRequest) (resp types.UpdateOrderResponse, err error) {
	defer utils.RecordTime(time.Now())()
	Price, _ := strconv.ParseInt(req.Price, 10, 64)
	GoodID, _ := strconv.ParseInt(req.GoodID, 10, 64)
	AddressID, _ := strconv.ParseInt(req.AddressID, 10, 64)
	Order := model.Order{
		TimeModel: model.TimeModel{},
		Price:     Price,
		GoodID:    GoodID,
		AddressID: AddressID,
	}

	err = repository.NewShopRequest(global.DB).UpdateOrder(Order)
	if err != nil {
		zlog.CtxErrorf(ctx, "订单更新到数据库失败: %v", err)
		return types.UpdateOrderResponse{}, errors.New("订单信息更新失败")
	}
	zlog.CtxInfof(ctx, "订单更新数据库处理完成: %+v", Order)
	return resp, nil
}
func (l *ShopLogic) SearchOrder(ctx context.Context, req types.SearchOrderRequest) (resp types.SearchOrderResponse, err error) {
	defer utils.RecordTime(time.Now())()
	Response, err := repository.NewShopRequest(global.DB).SearchOrder(req.KeyWord, req.OrderBy)
	if err != nil {
		zlog.CtxErrorf(ctx, "搜索商品失败: %v", err)
		return types.SearchOrderResponse{}, errors.New("搜索商品失败")
	}
	zlog.CtxInfof(ctx, "商品搜索完成: %+v", Response)
	resp.Orders = make([]types.Order, 0, len(Response))
	for _, Order := range Response {
		resp.Orders = append(resp.Orders, types.Order{
			ID:          strconv.FormatInt(Order.ID, 10),
			CreatedTime: strconv.FormatInt(Order.CreatedTime, 10),
			UpdatedTime: strconv.FormatInt(Order.UpdatedTime, 10),
			OwnerID:     strconv.FormatInt(Order.OwnerID, 10),
			IsDelete:    Order.IsDelete,
			Price:       strconv.FormatInt(Order.Price, 10),
			Status:      strconv.Itoa(Order.Status),
			Pay:         Order.Pay,
			PayBy:       Order.PayBy,
			PayTime:     strconv.FormatInt(Order.PayTime, 10),
		})
	}
	return resp, nil
}

func (l *ShopLogic) GetGoodInfo(ctx context.Context, req types.GetGoodInfoRequest) (resp types.GetGoodInfoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// 从数据库获取订单列表
	GoodID, _ := strconv.ParseInt(req.GoodID, 10, 64)
	Good, err := repository.NewShopRequest(global.DB).GetGoodInfo(GoodID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频列表失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	resp.Good = types.Good{
		ID:          strconv.FormatInt(Good.ID, 10),
		CreatedTime: strconv.FormatInt(Good.CreatedTime, 10),
		UpdatedTime: strconv.FormatInt(Good.UpdatedTime, 10),
		Name:        Good.Name,
		Info:        Good.Info,
		Category:    Good.Category,
		Tags:        Good.Tags,
		Img:         Good.Img,
		OriginPrice: Good.OriginPrice,
		SellPrice:   Good.SellPrice,
		Stock:       strconv.FormatInt(Good.Stock, 10),
		Status:      strconv.Itoa(Good.Status),
		OwnerID:     strconv.FormatInt(Good.OwnerID, 10),
	}
	return resp, nil
}
func (l *ShopLogic) CreateGood(ctx context.Context, req types.CreateGoodRequest) (resp types.CreateGoodResponse, err error) {
	defer utils.RecordTime(time.Now())()
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	Stock, _ := strconv.ParseInt(req.Stock, 10, 64)

	GoodID := global.SnowflakeNode.Generate().Int64()
	Good := model.Good{
		ID:          GoodID,
		TimeModel:   model.TimeModel{},
		Name:        req.Name,
		Info:        req.Info,
		Category:    req.Category,
		Tags:        req.Tags,
		Img:         req.Img,
		OriginPrice: req.OriginPrice,
		SellPrice:   req.SellPrice,
		Stock:       Stock,
		Status:      0,
		OwnerID:     OwnerID,
	}

	err = repository.NewShopRequest(global.DB).CreateGood(Good)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存商品到数据库失败: %v", err)
		return types.CreateGoodResponse{}, errors.New("保存商品信息失败")
	}
	zlog.CtxInfof(ctx, "商品信息保存到数据库成功, ID: %d", GoodID)

	zlog.CtxInfof(ctx, "商品保存数据库处理完成: %+v", Good)
	return resp, nil
}
func (l *ShopLogic) UpdateGood(ctx context.Context, req types.UpdateGoodRequest) (resp types.UpdateGoodResponse, err error) {
	defer utils.RecordTime(time.Now())()
	Stock, _ := strconv.ParseInt(req.Stock, 10, 64)

	GoodID := global.SnowflakeNode.Generate().Int64()
	Good := model.Good{
		ID:          GoodID,
		TimeModel:   model.TimeModel{},
		Name:        req.Name,
		Info:        req.Info,
		Category:    req.Category,
		CategoryId:  0,
		Tags:        req.Tags,
		Img:         req.Img,
		OriginPrice: req.OriginPrice,
		SellPrice:   req.SellPrice,
		Stock:       Stock,
		Status:      0,
	}

	err = repository.NewShopRequest(global.DB).UpdateGood(Good)
	if err != nil {
		zlog.CtxErrorf(ctx, "更新商品到数据库失败: %v", err)
		return types.UpdateGoodResponse{}, errors.New("更新商品信息失败")
	}
	zlog.CtxInfof(ctx, "商品更新到数据库成功, ID: %d", GoodID)

	zlog.CtxInfof(ctx, "商品更新数据库处理完成: %+v", Good)
	return resp, nil
}
func (l *ShopLogic) SearchGood(ctx context.Context, req types.SearchGoodRequest) (resp types.SearchGoodResponse, err error) {
	defer utils.RecordTime(time.Now())()
	GoodsCategoryId, _ := strconv.ParseInt(req.GoodsCategoryId, 10, 64)
	Response, err := repository.NewShopRequest(global.DB).SearchGood(GoodsCategoryId, req.KeyWord, req.OrderBy)
	if err != nil {
		zlog.CtxErrorf(ctx, "搜索商品失败: %v", err)
		return types.SearchGoodResponse{}, errors.New("搜索商品失败")
	}
	zlog.CtxInfof(ctx, "商品搜索完成: %+v", Response)
	resp.Goods = make([]types.Good, 0, len(Response))
	for _, Good := range Response {
		resp.Goods = append(resp.Goods, types.Good{
			ID:          strconv.FormatInt(Good.ID, 10),
			CreatedTime: strconv.FormatInt(Good.CreatedTime, 10),
			UpdatedTime: strconv.FormatInt(Good.UpdatedTime, 10),
			Name:        Good.Name,
			Info:        Good.Info,
			Category:    Good.Category,
			Tags:        Good.Tags,
			Img:         Good.Img,
			OriginPrice: Good.OriginPrice,
			SellPrice:   Good.SellPrice,
			Stock:       strconv.FormatInt(Good.Stock, 10),
			Status:      strconv.Itoa(Good.Status),
			OwnerID:     strconv.FormatInt(Good.OwnerID, 10),
		})
	}
	return resp, nil
}

func (l *ShopLogic) GetCartInfo(ctx context.Context, req types.GetCartInfoRequest) (resp types.GetCartInfoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// 从数据库获取订单列表
	CartID, _ := strconv.ParseInt(req.CartID, 10, 64)
	Cart, err := repository.NewShopRequest(global.DB).GetCartInfo(CartID)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频列表失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	resp.Cart = types.Cart{
		ID:          strconv.FormatInt(Cart.ID, 10),
		CreatedTime: strconv.FormatInt(Cart.CreatedTime, 10),
		UpdatedTime: strconv.FormatInt(Cart.UpdatedTime, 10),
		OwnerID:     strconv.FormatInt(Cart.OwnerID, 10),
		GoodsID:     strconv.FormatInt(Cart.GoodsID, 10),
		GoodsNum:    strconv.FormatInt(Cart.GoodsNum, 10),
		IsDelete:    Cart.IsDelete,
	}
	return resp, nil
}
func (l *ShopLogic) CreateCart(ctx context.Context, req types.CreateCartRequest) (resp types.CreateCartResponse, err error) {
	defer utils.RecordTime(time.Now())()
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)

	CartID := global.SnowflakeNode.Generate().Int64()
	Cart := model.Cart{
		ID:        CartID,
		TimeModel: model.TimeModel{},
		OwnerID:   OwnerID,
		IsDelete:  false,
	}

	err = repository.NewShopRequest(global.DB).CreateCart(Cart)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存购物车到数据库失败: %v", err)
		return types.CreateCartResponse{}, errors.New("保存购物车信息失败")
	}
	zlog.CtxInfof(ctx, "购物车信息保存到数据库成功, ID: %d", CartID)

	zlog.CtxInfof(ctx, "购物车保存数据库处理完成: %+v", Cart)
	return resp, nil
}
