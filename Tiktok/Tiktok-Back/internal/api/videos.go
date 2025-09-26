package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
)

func GetVideos(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetVideosRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取视频请求: %v", req)
	resp, err := logic.NewVideosLogic().GetVideos(ctx, req)
	response.Response(c, resp, err)
}

// 获取视频点赞数量
func GetVideoLikes(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetVideoLikesRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取视频点赞数量请求: %v", req)
	resp, err := logic.NewVideosLogic().GetVideoLikes(ctx, req)
	response.Response(c, resp, err)
}

// 获取评论点赞数量
func GetCommentLikes(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetCommentLikesRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取视频点赞数量请求: %v", req)
	resp, err := logic.NewVideosLogic().GetCommentLikes(ctx, req)
	response.Response(c, resp, err)
}

// 获取评论
func GetComments(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetCommentsRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "评论视频请求: %v", req)
	resp, err := logic.NewVideosLogic().GetComments(ctx, req)
	response.Response(c, resp, err)
}

// 获取评论详情
func GetCommentAll(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetCommentAllRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "评论视频请求: %v", req)
	resp, err := logic.NewVideosLogic().GetCommentAll(ctx, req)
	response.Response(c, resp, err)
}

// 获取评论
func GetCommentsMember(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetCommentsMemberRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取评论数请求: %v", req)
	resp, err := logic.NewVideosLogic().GetCommentsMember(ctx, req)
	response.Response(c, resp, err)
}

// 保存视频
func CreateVideo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateVideoRequest](c)
	if err != nil {
		return
	}
	req.UserID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "保存视频请求: %v", req)
	resp, err := logic.NewVideosLogic().CreateVideo(ctx, req)
	response.Response(c, resp, err)
}

// LikeVideo点赞视频
func LikeVideo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.LikeVideoRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "点赞视频请求: %v", req)
	resp, err := logic.NewVideosLogic().LikeVideo(ctx, req)
	response.Response(c, resp, err)
}

// 点赞评论
func LikeComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.LikeCommentRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "评论视频请求: %v", req)
	resp, err := logic.NewVideosLogic().LikeComment(ctx, req)
	response.Response(c, resp, err)
}

// 评论视频
func CommentVideo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CommentVideoRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "评论视频请求: %v", req)
	resp, err := logic.NewVideosLogic().CommentVideo(ctx, req)
	response.Response(c, resp, err)
}

// 评论评论
func CommentComment(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CommentCommentRequest](c)
	if err != nil {
		return
	}
	req.OwnerID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "评论视频请求: %v", req)
	resp, err := logic.NewVideosLogic().CommentComment(ctx, req)
	response.Response(c, resp, err)
}
