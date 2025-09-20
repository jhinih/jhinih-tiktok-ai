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

type VideosLogic struct {
}

func NewVideosLogic() *VideosLogic {
	return &VideosLogic{}
}

func (l *VideosLogic) GetVideos(ctx context.Context, req types.GetVideosRequest) (resp types.GetVideosResponse, err error) {
	defer utils.RecordTime(time.Now())()
	Page, _ := strconv.ParseInt(req.Page, 10, 64)
	PageSize, _ := strconv.ParseInt(req.Page, 10, 64)

	// 记录请求参数
	zlog.CtxInfof(ctx, "获取视频列表请求参数: page=%d, pageSize=%d", Page, PageSize)
	// 设置默认排序方式
	if req.OrderBy == "" {
		req.OrderBy = "random"
	}
	zlog.CtxInfof(ctx, "使用排序方式: %s", req.OrderBy)

	// 从数据库获取视频列表
	videos, err := repository.NewVideosRequest(global.DB).GetVideos(Page, PageSize, req.OrderBy)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频列表失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	// 获取总记录数
	total, err := repository.NewVideosRequest(global.DB).GetTotalVideosCount()
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频总数失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	// 计算是否有更多数据
	hasMore := int64(Page*PageSize) < total

	// 构建响应数据
	resp.Data = make([]types.Video, 0, len(videos))
	for _, video := range videos {
		resp.Data = append(resp.Data, types.Video{
			ID:          strconv.FormatInt(video.ID, 10),
			CreatedTime: strconv.FormatInt(time.Now().Unix(), 10),
			UpdatedTime: strconv.FormatInt(time.Now().Unix(), 10),
			Title:       video.Title,
			Description: video.Description,
			URL:         video.URL,
			Cover:       video.Cover,
			Likes:       strconv.FormatInt(video.Likes, 10),
			Comments:    strconv.FormatInt(video.Comments, 10),
			Shares:      strconv.FormatInt(video.Shares, 10),
			UserID:      strconv.FormatInt(video.UserID, 10),
			PublishTime: video.PublishTime,
			Type:        video.Type,
			IsPrivate:   video.IsPrivate,
		})
	}

	// 设置分页元数据
	resp.Page = strconv.FormatInt(Page, 10)
	resp.PageSize = strconv.FormatInt(PageSize, 10)
	resp.Total = strconv.FormatInt(total, 10)
	resp.HasMore = hasMore

	// 记录响应数据
	zlog.CtxInfof(ctx, "返回视频列表响应: 当前页%s, 每页%s条, 共%s条, 是否有更多:%v",
		resp.Page, resp.PageSize, resp.Total, resp.HasMore)

	return resp, nil
}
func (l *VideosLogic) GetVideoLikes(ctx context.Context, req types.GetVideoLikesRequest) (resp types.GetVideoLikesResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始获取视频点赞数量, 视频: %s", req.VideoID)
	VideoID, _ := strconv.ParseInt(req.VideoID, 10, 64)
	var VideoLikes int64
	VideoLikes, err = repository.NewVideosRequest(global.DB).GetVideoLikes(VideoID)
	resp.VideoLikes = strconv.FormatInt(VideoLikes, 10)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频点赞数量失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}
func (l *VideosLogic) GetCommentLikes(ctx context.Context, req types.GetCommentLikesRequest) (resp types.GetCommentLikesResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始获取评论点赞数量, 评论: %s", req.CommentID)
	CommentID, _ := strconv.ParseInt(req.CommentID, 10, 64)
	var CommentLikes int64
	CommentLikes, err = repository.NewVideosRequest(global.DB).GetCommentLikes(CommentID)
	resp.CommentLikes = strconv.FormatInt(CommentLikes, 10)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取评论点赞数量失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}
func (l *VideosLogic) GetComments(ctx context.Context, req types.GetCommentsRequest) (resp types.GetCommentsResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// ID 转化为 int64
	ID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	beforeID := global.SnowflakeNode.Generate().Int64()
	Page, _ := strconv.ParseInt(req.Page, 10, 64)
	PageSize, _ := strconv.ParseInt(req.Page, 10, 64)
	// 设置默认排序方式
	if req.OrderBy == "" {
		req.OrderBy = "latest"
	}
	//beforeID, err = strconv.ParseInt(beforeID, 10, 64)
	//if err != nil {
	//	zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.BeforeID, err)
	//	return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	//}
	// 从数据库中查询评论
	var comments []model.Comment
	if req.IsVideo {
		comments, err = repository.NewVideosRequest(global.DB).GetVideoComments(ID, beforeID, Page, PageSize, req.OrderBy)
	} else {
		comments, err = repository.NewVideosRequest(global.DB).GetCommentComments(ID, beforeID, Page, PageSize, req.OrderBy)
	}
	if err != nil {
		zlog.CtxErrorf(ctx, "查询评论失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	for _, comment := range comments {
		resp.Comments = append(resp.Comments, types.Comment{
			ID:          strconv.FormatInt(comment.ID, 10),
			Content:     comment.Content,
			VideoID:     strconv.FormatInt(comment.VideoID, 10),
			FatherID:    strconv.FormatInt(comment.FatherID, 10),
			UserID:      strconv.FormatInt(comment.UserID, 10),
			OwnerId:     strconv.FormatInt(comment.OwnerId, 10),
			Likes:       strconv.FormatInt(comment.Likes, 10),
			Comments:    strconv.FormatInt(comment.Comments, 10),
			CreatedTime: strconv.FormatInt(comment.CreatedTime, 10),
			UpdatedTime: strconv.FormatInt(comment.UpdatedTime, 10),
		})
	}
	resp.Length = strconv.Itoa(len(resp.Comments))
	return resp, err
}
func (l *VideosLogic) GetCommentAll(ctx context.Context, req types.GetCommentAllRequest) (resp types.GetCommentAllResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// ID 转化为 int64
	ID, err := strconv.ParseInt(req.CommentID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.CommentID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	// 从数据库中查询评论
	var comment model.Comment
	comment, err = repository.NewVideosRequest(global.DB).GetCommentAll(ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "查询评论详情失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	Comment := types.Comment{
		ID:          strconv.FormatInt(comment.ID, 10),
		Content:     comment.Content,
		VideoID:     strconv.FormatInt(comment.VideoID, 10),
		FatherID:    strconv.FormatInt(comment.FatherID, 10),
		UserID:      strconv.FormatInt(comment.UserID, 10),
		OwnerId:     strconv.FormatInt(comment.OwnerId, 10),
		Likes:       strconv.FormatInt(comment.Likes, 10),
		Comments:    strconv.FormatInt(comment.Comments, 10),
		CreatedTime: strconv.FormatInt(comment.CreatedTime, 10),
		UpdatedTime: strconv.FormatInt(comment.UpdatedTime, 10),
	}
	resp.Comment = Comment
	return resp, nil
}
func (l *VideosLogic) GetCommentsMember(ctx context.Context, req types.GetCommentsMemberRequest) (resp types.GetCommentsMemberResponse, err error) {
	defer utils.RecordTime(time.Now())()
	// ID 转化为 int64
	ID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "%v 转换 int64 错误: %v", req.ID, err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	// 从数据库中查询评论
	var comments int64
	if req.IsVideo {
		comments, err = repository.NewVideosRequest(global.DB).GetVideoCommentsMember(ID)
	} else {
		comments, err = repository.NewVideosRequest(global.DB).GetCommentCommentsMember(ID)
	}
	if err != nil {
		zlog.CtxErrorf(ctx, "查询评论数失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	resp.Member = strconv.FormatInt(comments, 10)
	return resp, nil
}

func (l *VideosLogic) CreateVideo(ctx context.Context, req types.CreateVideoRequest) (resp types.CreateVideoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	VideoID := global.SnowflakeNode.Generate().Int64()
	video := model.Video{
		ID: VideoID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
		UserID:      UserID,
		URL:         req.VideoPath,
		Cover:       req.CoverPath,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		Type:        req.Type,
		PublishTime: time.Now(),
		Likes:       0,
		Comments:    0,
		Shares:      0,
	}
	err = repository.NewVideosRequest(global.DB).CreateVideo(video)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存视频信息到数据库失败: %v, 视频路径: %s", err, req.VideoPath)
		return types.CreateVideoResponse{}, errors.New("保存视频信息失败")
	}
	zlog.CtxInfof(ctx, "视频信息保存到数据库成功, ID: %d", VideoID)

	zlog.CtxInfof(ctx, "视频保存数据库处理完成: %+v", video)
	return resp, nil
}

func (l *VideosLogic) LikeVideo(ctx context.Context, req types.LikeVideoRequest) (resp types.LikeVideoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始处理视频点赞, 视频: %s, 作者ID: %s,用户ID: %s", req.VideoID, req.OwnerID, req.UserID)
	VideoID, _ := strconv.ParseInt(req.VideoID, 10, 64)
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	videolike := model.VideoLike{
		VideoID: VideoID,
		OwnerID: OwnerID,
		UserID:  UserID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
	err = repository.NewVideosRequest(global.DB).LikeVideo(videolike)
	if err != nil {
		zlog.CtxErrorf(ctx, "视频点赞失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}
func (l *VideosLogic) LikeComment(ctx context.Context, req types.LikeCommentRequest) (resp types.LikeCommentResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始处理评论点赞, 评论: %s, 作者ID: %s,用户ID: %s", req.CommentID, req.OwnerID, req.UserID)
	CommentID, _ := strconv.ParseInt(req.CommentID, 10, 64)
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	commentlike := model.CommentLike{
		CommentID: CommentID,
		OwnerID:   OwnerID,
		UserID:    UserID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
	err = repository.NewVideosRequest(global.DB).LikeComment(commentlike)
	if err != nil {
		zlog.CtxErrorf(ctx, "评论点赞失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}

func (l *VideosLogic) CommentVideo(ctx context.Context, req types.CommentVideoRequest) (resp types.CommentVideoResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始评论视频, 视频: %s,作者:%s,用户:%s,评论:%s", req.VideoID, req.OwnerID, req.UserID, req.Content)
	VideoID, _ := strconv.ParseInt(req.VideoID, 10, 64)
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	ID := global.SnowflakeNode.Generate().Int64()
	VideoComment := model.Comment{
		ID:       ID,
		VideoID:  VideoID,
		UserID:   UserID,
		Content:  req.Content,
		OwnerId:  OwnerID,
		FatherID: ID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
		Likes:    0,
		Comments: 0,
	}
	err = repository.NewVideosRequest(global.DB).CommentVideo(VideoComment)
	if err != nil {
		zlog.CtxErrorf(ctx, "评论视频失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}
func (l *VideosLogic) CommentComment(ctx context.Context, req types.CommentCommentRequest) (resp types.CommentCommentResponse, err error) {
	defer utils.RecordTime(time.Now())()
	zlog.CtxInfof(ctx, "开始评论评论, 评论: %s,作者:%s,用户:%s,评论:%s", req.FatherID, req.OwnerID, req.UserID, req.Content)
	OwnerID, _ := strconv.ParseInt(req.OwnerID, 10, 64)
	UserID, _ := strconv.ParseInt(req.UserID, 10, 64)
	FatherID, _ := strconv.ParseInt(req.FatherID, 10, 64) //string---int64
	ID := global.SnowflakeNode.Generate().Int64()
	VideoComment := model.Comment{
		ID:       ID,
		UserID:   UserID,
		Content:  req.Content,
		FatherID: FatherID,
		OwnerId:  OwnerID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
		Likes:    0,
		Comments: 0,
	}
	err = repository.NewVideosRequest(global.DB).CommentComment(VideoComment)
	if err != nil {
		zlog.CtxErrorf(ctx, "评论评论失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}
	return resp, nil
}
