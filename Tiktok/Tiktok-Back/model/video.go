package model

import "time"

// Video 视频数据结构
type Video struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Cover       string `json:"cover"`
	Likes       int64  `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	Comments    int64  `json:"comments" gorm:"not null;column:comments;type:int;comment:评论数"`
	Shares      int64  `json:"shares" gorm:"not null;column:shares;type:int;comment:分享数"`
	UserID      int64  `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	PublishTime time.Time
	Type        string `json:"type" gorm:"column:type;type:varchar(63);comment:类型"`
	IsPrivate   bool   `json:"is_private" gorm:"not null;column:is_private;type:bool;comment:是否私密"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoLike struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	VideoID int64 `json:"video_id" gorm:"column:video_id;type:bigint;comment:视频ID;uniqueIndex:idx_video_user_unique"`
	UserID  int64 `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID;uniqueIndex:idx_video_user_unique"`
	OwnerID int64 `json:"owner_id" gorm:"column:owner_id;type:bigint;comment:作者ID;uniqueIndex:idx_video_user_unique"`
}

func (VideoLike) TableName() string {
	return "video_like"
}

type Comment struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Content  any   `json:"content" gorm:"column:content;type:text;comment:内容"`
	VideoID  int64 `json:"video_id" gorm:"column:video_id;type:bigint;comment:视频ID;index"`
	FatherID int64 `json:"father_id" gorm:"column:father_id;type:bigint;comment:父评论ID"`
	UserID   int64 `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	OwnerId  int64 `json:"owner_id" gorm:"column:owner_id;type:bigint;comment:作者ID"`
	Likes    int64 `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	Comments int64 `json:"comments" gorm:"not null;column:comments;type:int;comment:评论数"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentLike struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	CommentID int64 `json:"comment_id" gorm:"column:comment_id;type:bigint;comment:评论ID;uniqueIndex:idx_comment_user_unique"`
	UserID    int64 `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID;uniqueIndex:idx_comment_user_unique"`
	OwnerID   int64 `json:"owner_id" gorm:"column:owner_id;type:bigint;comment:作者ID"`
}

func (CommentLike) TableName() string {
	return "comment_like"
}
