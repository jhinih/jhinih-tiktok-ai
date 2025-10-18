package model

type Good struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Name       string `json:"name" gorm:"column:name"`
	Info       string `json:"info" gorm:"column:info"`
	Category   string `json:"category" gorm:"column:category"`
	CategoryId int64  `json:"category_id" gorm:"column:category_id"`
	Tags       string `json:"tags" gorm:"column:tags"`
	Img        string `json:"img" gorm:"column:img"`

	OriginPrice string `json:"originPrice" gorm:"column:origin_price"`
	SellPrice   string `json:"sellPrice" gorm:"column:sell_price"`

	Stock  int64 `json:"stock" gorm:"column:stock"`
	Status int   `json:"status" gorm:"column:status;type:tinyint;type:tinyint"`

	OwnerID int64 `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
}

func (Good) TableName() string {
	return "Good"
}
