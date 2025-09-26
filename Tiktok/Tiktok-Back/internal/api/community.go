package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
)

func CreateCommunity(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateCommunityRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "创建群聊请求: %v", req)
	resp, err := logic.NewCommunityLogic().CreateCommunity(ctx, req)
	response.Response(c, resp, err)
}
func JoinCommunity(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.JoinCommunityRequest](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "加入群聊请求: %v", req)
	resp, err := logic.NewCommunityLogic().JoinCommunity(ctx, req)
	response.Response(c, resp, err)
}
func LoadCommunity(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.LoadCommunityRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "加载群聊请求: %v", req)
	resp, err := logic.NewCommunityLogic().LoadCommunity(ctx, req)
	response.Response(c, resp, err)
}
