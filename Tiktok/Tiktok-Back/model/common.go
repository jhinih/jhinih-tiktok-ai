package model

import (
	"gorm.io/gorm"
	"time"
)

// CommonModel 每张表都有的四个东西，最好不要用 gorm.model（虽然他们一模一样）
type CommonModel struct {
	ID        int64 `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TimeModel struct {
	CreatedTime int64 `gorm:"column:created_time;type:bigint"`
	UpdatedTime int64 `gorm:"column:updated_time;type:bigint"`
}

//func (b *CommonModel) BeforeCreate(db *gorm.DB) error {
//	// 生成雪花ID
//	if b.ID == 0 {
//		b.ID = snowflake.GetIntId(global.Node)
//	}
//
//	return nil
//}

func (b *TimeModel) BeforeCreate(db *gorm.DB) error {
	// 生成雪花ID
	b.CreatedTime = time.Now().UnixMilli()
	b.UpdatedTime = time.Now().UnixMilli()
	return nil
}
