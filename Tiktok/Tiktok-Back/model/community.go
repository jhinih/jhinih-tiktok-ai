package model

type Community struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Name      string `json:"name" gorm:"column:name;type:varchar(255);size:255;unique;not null;"`
	OwnerID   int64  `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	OwnerName string `json:"owner_name" gorm:"column:owner_name;type:varchar(255);size:255;"`
	Img       string `json:"img" gorm:"column:img;type:varchar(255);size:255;"`
	Desc      string `json:"desc" gorm:"column:desc;type:varchar(255);size:255;"`
}

func (Community) TableName() string {
	return "Community"
}
