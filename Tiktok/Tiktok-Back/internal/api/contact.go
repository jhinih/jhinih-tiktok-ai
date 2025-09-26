package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
)

// 加好友
func AddFriend(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AddFriendRequest](c)
	if err != nil {
		return
	}
	// 从token中获取用户ID
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "加好友请求: %v", req)
	resp, err := logic.NewContactLogic().AddFriend(ctx, req)
	response.Response(c, resp, err)
}

// 好友列表
func GetFriendList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetFriendListRequest](c)
	if err != nil {
		return
	}
	// 从token中获取用户ID
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "获取好友列表请求: %v", req)
	resp, err := logic.NewContactLogic().GetFriendList(c, req)
	response.Response(c, resp, err)
}

// 在线列表
func GetUserListOnline(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetUserListOnlineRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取在线列表请求: %v", req)
	resp, err := logic.NewContactLogic().GetUserListOnline(c)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取在线列表失败: %v", err)
		response.Response(c, resp, err)
		return
	}
	response.Response(c, resp, nil)
}

// 群友列表
func GetGroupUsers(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetGroupUsersRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "搜索群友请求: %v", req)
	resp, err := logic.NewContactLogic().GetGroupUsers(c, req)
	response.Response(c, resp, err)
}

// 群列表
func GetGroupList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetGroupListRequest](c)
	if err != nil {
		return
	}
	// 从token中获取用户ID
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "搜索群列表请求: %v", req)
	resp, err := logic.NewContactLogic().GetGroupList(c, req)
	response.Response(c, resp, err)
}
