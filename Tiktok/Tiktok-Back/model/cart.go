package model

type Cart struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	OwnerID  int64  `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	GoodsID  int64  `json:"goods_id" gorm:"column:goods_id;type:bigint"`
	GoodsNum int64  `json:"goods_num" gorm:"column:goods_num;type:bigint"`
	Goods    []Good `json:"goods" gorm:"foreignKey:GoodsID;references:ID"`

	IsDelete bool `json:"is_delete" gorm:"column:is_delete;type:tinyint;type:tinyint"`
}

func (Cart) TableName() string {
	return "Cart"
}
