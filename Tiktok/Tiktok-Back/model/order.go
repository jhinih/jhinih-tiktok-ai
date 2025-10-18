package model

type Order struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	OrderID int64 `json:"order_id" gorm:"column:order_id"`
	OwnerID int64 `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`

	GoodID    int64 `json:"good_id" gorm:"column:good_id;type:bigint;type:bigint"`
	AddressID int64 `json:"address_id" gorm:"column:address_id;type:bigint;type:bigint"`
	Price     int64 `json:"price" gorm:"column:price;type:bigint;type:bigint"`

	Status int `json:"status" gorm:"column:status;type:tinyint;type:tinyint"`

	Pay     bool   `json:"pay" gorm:"column:pay;type:tinyint;type:tinyint"`
	PayBy   string `json:"pay_by" gorm:"column:pay_by;type:varchar(50);type:varchar(50)"`
	PayTime int64  `json:"pay_time" gorm:"column:pay_time;type:bigint"`

	IsDelete bool `json:"is_delete" gorm:"column:is_delete;type:tinyint;type:tinyint"`
}

func (Order) TableName() string {
	return "Order"
}
