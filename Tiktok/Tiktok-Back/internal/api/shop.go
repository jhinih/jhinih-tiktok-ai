package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
)

// GetOrderInfo 获取订单信息
func GetOrderInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetOrderInfoRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "获取订单信息请求: %v", req)
	resp, err := logic.NewShopLogic().GetOrderInfo(ctx, req)
	response.Response(c, resp, err)
}
func CreateOrder(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateOrderRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "创建订单请求: %v", req)
	resp, err := logic.NewShopLogic().CreateOrder(ctx, req)
	response.Response(c, resp, err)
}

// GetGoodInfo 获取商品信息
func GetGoodInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetGoodInfoRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "获取商品信息请求: %v", req)
	resp, err := logic.NewShopLogic().GetGoodInfo(ctx, req)
	response.Response(c, resp, err)
}

func CreateGood(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateGoodRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "创建商品请求: %v", req)
	resp, err := logic.NewShopLogic().CreateGood(ctx, req)
	response.Response(c, resp, err)
}

func SearchGood(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SearchGoodRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "创建商品请求: %v", req)
	resp, err := logic.NewShopLogic().SearchGood(ctx, req)
	response.Response(c, resp, err)
}

// GetCartInfo 获取商品信息
func GetCartInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetCartInfoRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "获取购物车信息请求: %v", req)
	resp, err := logic.NewShopLogic().GetCartInfo(ctx, req)
	response.Response(c, resp, err)
}

func CreateCart(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateCartRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "创建购物车请求: %v", req)
	resp, err := logic.NewShopLogic().CreateCart(ctx, req)
	response.Response(c, resp, err)
}
