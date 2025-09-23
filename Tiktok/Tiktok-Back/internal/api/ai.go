package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"github.com/gin-gonic/gin"
)

func CommonAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "普通AI请求: %v", req)
	resp, err := logic.NewAILogic().CommonAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}
func VideoAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "普通AI请求: %v", req)
	resp, err := logic.NewAILogic().CommonAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}
func SendCodeAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "发送邮件AI请求: %v", req)
	resp, err := logic.NewAILogic().SendCodeAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}

func GetUserInfoAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "获取用户信息AI请求: %v", req)
	resp, err := logic.NewAILogic().GetUserInfoAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}

func GetVideoAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "获取视频AI请求: %v", req)
	resp, err := logic.NewAILogic().GetVideoAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}

func AllAI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "AI请求: %v", req)
	resp, err := logic.NewAILogic().AllAI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}
func AI(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AIRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}

	zlog.CtxInfof(ctx, "AI请求: %v", req)
	resp, err := logic.NewAILogic().AI(c, req)
	if err != nil {
		response.Response(c, nil, err)
	}
	response.Response(c, resp, nil)
}
