package types

import (
	"time"
)

// 获取视频列表请求
type GetVideosRequest struct {
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
	OrderBy  string `json:"order_by"` // random:随机, latest:最新, popular:最热
}

type Video struct {
	ID          string `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Cover       string `json:"cover"`
	Likes       string `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	Comments    string `json:"comments" gorm:"not null;column:comments;type:int;comment:评论数"`
	Shares      string `json:"shares" gorm:"not null;column:shares;type:int;comment:分享数"`
	UserID      string `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	PublishTime time.Time
	Type        string `json:"type" gorm:"column:type;type:varchar(63);comment:类型"`
	IsPrivate   bool   `json:"is_private" gorm:"not null;column:is_private;type:bool;comment:是否私密"`
}

// 获取视频列表响应
type GetVideosResponse struct {
	Data     []Video `json:"data"`
	Page     string  `json:"page"`
	PageSize string  `json:"page_size"`
	Total    string  `json:"total"`
	HasMore  bool    `json:"has_more"`
}

// 获取视频点赞数请求
type GetVideoLikesRequest struct {
	VideoID string `form:"video_id" binding:"required"`
}

// 获取视频点赞数响应
type GetVideoLikesResponse struct {
	VideoLikes string `json:"video_likes"`
}

// 获取评论点赞数请求
type GetCommentLikesRequest struct {
	CommentID string `form:"comment_id" binding:"required"`
}

// 获取评论点赞数响应
type GetCommentLikesResponse struct {
	CommentLikes string `json:"comment_likes"`
}

// 获取评论评论请求
type GetCommentsRequest struct {
	ID      string `form:"id"`
	IsVideo bool   `form:"is_video"`
	//BeforeID string `form:"before_id"`
	Page     string `form:"page"`
	PageSize string `form:"page_size"`
	OrderBy  string `form:"order_by"` // random:随机, latest:最新, popular:最热
}

type Comment struct {
	ID          string `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	Content     any    `json:"content" gorm:"column:content;type:text;comment:内容"`
	VideoID     string `json:"video_id" gorm:"column:video_id;type:bigint;comment:视频ID;index"`
	FatherID    string `json:"father_id" gorm:"column:father_id;type:bigint;comment:父评论ID"`
	UserID      string `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	OwnerId     string `json:"owner_id" gorm:"column:owner_id;type:bigint;comment:作者ID"`
	Likes       string `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	Comments    string `json:"comments" gorm:"not null;column:comments;type:int;comment:评论数"`
}

// 获取评论评论响应
type GetCommentsResponse struct {
	Comments []Comment `json:"comments"`
	Length   string    `json:"length"`
}

// 获取评论请求
type GetCommentAllRequest struct {
	CommentID string `form:"comment_id" binding:"required"`
}

// 获取评论响应
type GetCommentAllResponse struct {
	Comment Comment `json:"comment"`
}

type GetCommentsMemberRequest struct {
	ID      string `form:"id"`
	IsVideo bool   `form:"is_video"`
}

// 获取评论评论响应
type GetCommentsMemberResponse struct {
	Member string `json:"member"`
}

// 上传视频请求参数
type CreateVideoRequest struct {
	VideoPath   string `json:"video_path"`
	CoverPath   string `json:"cover_path"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	UserID      string `json:"user_id" binding:"-"`
	Type        string `json:"type"`
}

// 上传视频响应
type CreateVideoResponse struct {
}

// 点赞视频请求参数
type LikeVideoRequest struct {
	VideoID string `json:"video_id" binding:"required"`
	OwnerID string `json:"owner_id" binding:"-"`
	UserID  string `json:"user_id"`
}

// 点赞视频响应
type LikeVideoResponse struct {
}

type LikeCommentRequest struct {
	CommentID string `json:"comment_id" binding:"required"`
	OwnerID   string `json:"owner_id" binding:"-"`
	UserID    string `json:"user_id" `
}

// 点赞评论响应
type LikeCommentResponse struct {
}

// 评论视频请求
type CommentVideoRequest struct {
	VideoID string `json:"video_id" binding:"required"`
	Content any    `json:"content"`
	OwnerID string `json:"owner_id" binding:"-"`
	UserID  string `json:"user_id"`
}

// 评论视频响应
type CommentVideoResponse struct {
}

// 评论视频响应
type CommentCommentRequest struct {
	Content  any    `json:"content"`
	OwnerID  string `json:"owner_id" binding:"-"`
	UserID   string `json:"user_id"`
	FatherID string `json:"father_id"`
}

// 评论视频响应
type CommentCommentResponse struct {
}
