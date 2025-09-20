package model

type Contact struct {
	ID int64 `json:"id,string" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	OwnerID    int64  `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	OwnerName  string `json:"owner_name" gorm:"column:owner_name;type:varchar(255);size:255;"`
	TargetID   int64  `json:"target_id,string" gorm:"column:target_id;type:bigint;type:bigint"`
	TargetName string `json:"target_name" gorm:"column:target_name;type:varchar(255);size:255;"`
	Type       int    //对应的类型  1好友  2群  3xx
	Desc       string `json:"desc" gorm:"column:desc;type:varchar(255);size:255;"`
}

func (table *Contact) TableName() string {
	return "contact"
}
