package repository

import (
	"Tiktok/model"
	"Tiktok/repository/list"
	"errors"
	"gorm.io/gorm"
	"time"
)

type VideosRequest struct {
	DB *gorm.DB
}

func NewVideosRequest(db *gorm.DB) *VideosRequest {
	return &VideosRequest{
		DB: db,
	}
}

// 获取视频
func (r *VideosRequest) GetVideoByVideoID(videoId int64) (model.Video, error) {
	var tableVideo model.Video
	tableVideo.ID = videoId
	//Init()
	result := r.DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}

// 获取作者视频列表
func (r *VideosRequest) GetVideosByUserID(UserID int64) ([]model.Video, error) {
	//建立结果集接收
	var data []model.Video
	result := r.DB.Where(&model.Video{UserID: UserID}).Find(&data)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// 获取这个时间之前的一些视频
func (r *VideosRequest) GetVideosByLastTime(lastTime time.Time) ([]model.Video, error) {
	videos := make([]model.Video, 5)
	result := r.DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(5).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// 获取分页视频列表
func (r *VideosRequest) GetVideos(page, pageSize int64, orderBy string) ([]model.Video, error) {
	//// 设置默认分页值
	//if page <= 0 {
	//	page = 1
	//}
	//if pageSize <= 0 {
	//	pageSize = 10
	//}
	//
	//// 设置默认排序方式
	//if orderBy == "" {
	//	orderBy = "random"
	//}
	//
	//offset := (page - 1) * pageSize
	//var videos []model.Video
	//var subQuery *gorm.DB
	//
	//// 根据排序方式构建查询
	//switch orderBy {
	//case "random":
	//	// 获取所有符合条件的ID
	//	var allIDs []int64
	//	err := r.DB.Model(&model.Video{}).
	//		Select("id").
	//		Where("is_private = ?", false).
	//		Find(&allIDs).Error
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	// 随机排序ID
	//	rand.Shuffle(len(allIDs), func(i, j int) {
	//		allIDs[i], allIDs[j] = allIDs[j], allIDs[i]
	//	})
	//
	//	// 分页处理
	//	start := offset
	//	if start > len(allIDs) {
	//		start = len(allIDs)
	//	}
	//	end := start + pageSize
	//	if end > len(allIDs) {
	//		end = len(allIDs)
	//	}
	//	pageIDs := allIDs[start:end]
	//
	//	// 获取分页视频数据
	//	err = r.DB.Model(&model.Video{}).
	//		Where("id IN (?)", pageIDs).
	//		Find(&videos).Error
	//	return videos, err
	//case "latest":
	//	// 按创建时间降序
	//	subQuery = r.DB.Model(&model.Video{}).
	//		Select("id").
	//		Where("is_private = ?", false).
	//		Order("created_at DESC").
	//		Offset(offset).
	//		Limit(pageSize)
	//case "popular":
	//	// 按热度(点赞数)降序
	//	subQuery = r.DB.Model(&model.Video{}).
	//		Select("id").
	//		Where("is_private = ?", false).
	//		Order("likes DESC, created_at DESC").
	//		Offset(offset).
	//		Limit(pageSize)
	//default:
	//	// 默认随机排序
	//	subQuery = r.DB.Model(&model.Video{}).
	//		Select("id").
	//		Where("is_private = ?", false).
	//		Order("RAND()").
	//		Offset(offset).
	//		Limit(pageSize)
	//}
	//
	//// 获取分页视频数据
	//err := r.DB.Model(&model.Video{}).
	//	Where("id IN (?)", subQuery).
	//	Find(&videos).Error
	//
	//return videos, err
	var videos []model.Video
	videos, _, err := list.Query[model.Video](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("is_private = ?", false)
		},
		Order: orderBy,
		//Preloads: []string{"Author"},
		PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	})
	return videos, err
}

// 获取视频总数
func (r *VideosRequest) GetTotalVideosCount() (int64, error) {
	var count int64
	err := r.DB.Model(&model.Video{}).
		Where("is_private = ?", false).
		Count(&count).Error
	return count, err
}

// 获取点赞
func (r *VideosRequest) GetVideoLikes(VideoID int64) (int64, error) {
	var Video model.Video
	if err := r.DB.Select("likes").Where("id = ?", VideoID).First(&Video).Error; err != nil {
		return 0, err
	}
	return Video.Likes, nil
}
func (r *VideosRequest) GetCommentLikes(CommentID int64) (int64, error) {
	var Comment model.Comment
	if err := r.DB.Select("likes").Where("id = ?", CommentID).First(&Comment).Error; err != nil {
		return 0, err
	}
	return Comment.Likes, nil
}

// 获取评论
func (r *VideosRequest) GetVideoComments(VideoID, before, page, pageSize int64, orderBy string) ([]model.Comment, error) {
	//var comments []model.Comment
	//err := r.DB.Model(&model.Comment{}).Where("video_id = ? AND id < ? AND father_id = father_id", VideoID, before).Order("id DESC").Limit(count).Find(&comments).Error
	//return comments, err
	var comments []model.Comment
	comments, _, err := list.Query[model.Comment](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("video_id = ? AND id < ? AND father_id = father_id", VideoID, before)
		},
		Order: orderBy,
		//Preloads: []string{"Author"},
		PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	})
	return comments, err
}
func (r *VideosRequest) GetCommentComments(FatherID, before, page, pageSize int64, orderBy string) ([]model.Comment, error) {
	//var comments []model.Comment
	//err := r.DB.Model(&model.Comment{}).Where("father_id = ? AND id < ? ", FatherID, before).Order("id DESC").Limit(count).Find(&comments).Error
	//return comments, err
	var comments []model.Comment
	comments, _, err := list.Query[model.Comment](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("father_id = ? AND id < ? ", FatherID, before)
		},
		Order: orderBy,
		//Preloads: []string{"Author"},
		PageInfo: list.PageInfo{Page: page, Limit: pageSize},
	})
	return comments, err
}
func (r *VideosRequest) GetCommentAll(ID int64) (model.Comment, error) {
	var comments model.Comment
	err := r.DB.Model(&model.Comment{}).Where("id = ?", ID).Order("id DESC").Find(&comments).Error
	return comments, err
}

// 保存视频
func (r *VideosRequest) CreateVideo(video model.Video) error {
	return r.DB.Create(&video).Error
}

// 点赞
func (r *VideosRequest) LikeVideo(videolike model.VideoLike) error {
	// 检查是否已点赞
	var like model.VideoLike
	if err := r.DB.Where("video_id = ? AND user_id = ?", videolike.VideoID, videolike.UserID).First(&like).Error; err == nil {
		// 已点赞，执行取消点赞
		if err := r.DB.Delete(&like).Error; err != nil {
			return err
		}
		// 减少视频点赞数
		return r.DB.Model(&model.Video{}).Where("id = ?", videolike.VideoID).
			Update("likes", gorm.Expr("likes - 1")).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 未点赞，执行点赞
	newLike := model.VideoLike{
		VideoID: videolike.VideoID,
		UserID:  videolike.UserID,
		OwnerID: videolike.OwnerID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
		},
	}
	if err := r.DB.Create(&newLike).Error; err != nil {
		return err
	}
	// 增加视频点赞数
	return r.DB.Model(&model.Video{}).Where("id = ?", newLike.VideoID).
		Update("likes", gorm.Expr("likes + 1")).Error
}
func (r *VideosRequest) LikeComment(commentlike model.CommentLike) error {
	// 检查是否已点赞
	var like model.CommentLike
	if err := r.DB.Where("comment_id = ? AND user_id = ?", commentlike.CommentID, commentlike.UserID).First(&like).Error; err == nil {
		// 已点赞，执行取消点赞
		if err := r.DB.Delete(&like).Error; err != nil {
			return err
		}
		// 减少评论点赞数
		return r.DB.Model(&model.Comment{}).Where("id = ?", commentlike.CommentID).
			Update("likes", gorm.Expr("likes - 1")).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// 未点赞，执行点赞
	newLike := model.CommentLike{
		CommentID: commentlike.CommentID,
		UserID:    commentlike.UserID,
		OwnerID:   commentlike.OwnerID,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
	if err := r.DB.Create(&newLike).Error; err != nil {
		return err
	}
	// 增加评论点赞数
	return r.DB.Model(&model.Comment{}).Where("id = ?", newLike.CommentID).
		Update("likes", gorm.Expr("likes + 1")).Error
}

// 创建评论
func (r *VideosRequest) CommentVideo(comment model.Comment) error {
	err := r.DB.Create(&comment).Error
	if err != nil {
		return err
	}
	// 更新视频评论数
	return r.DB.Model(&model.Video{}).Where("id = ?", comment.VideoID).Update("comments", gorm.Expr("comments + ?", 1)).Error
}
func (r *VideosRequest) CommentComment(comment model.Comment) error {
	err := r.DB.Create(&comment).Error
	if err != nil {
		return err
	}
	// 更新评论评论数
	return r.DB.Model(&model.Comment{}).Where("id = ?", comment.FatherID).Update("comments", gorm.Expr("comments + ?", 1)).Error
}

// 获取评论

// 获取评论数
func (r *VideosRequest) GetVideoCommentsMember(VideoID int64) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Video{}).
		Where("father_id = ?", VideoID).
		Count(&count).Error
	return count, err
}
func (r *VideosRequest) GetCommentCommentsMember(FatherID int64) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Comment{}).
		Where("father_id = ?", FatherID).
		Count(&count).Error
	return count, err
}

//func (r *VideosRequest) GetCommentCommentsMember(FatherID int64) (int64, error) {
//	var comments model.Comment
//	err := r.DB.Model(&model.Comment{}).Where("father_id = ?", FatherID).Order("id DESC").Find(&comments).Error
//	return comments.Comments, err
//}
