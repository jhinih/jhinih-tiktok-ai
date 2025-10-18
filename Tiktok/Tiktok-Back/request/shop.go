package request

import (
	"Tiktok/model"
	"gorm.io/gorm"
)

type ShopRequest struct {
	DB *gorm.DB
}

func NewShopRequest(db *gorm.DB) *ShopRequest {
	return &ShopRequest{
		DB: db,
	}
}

// GetOrderInfo 获取订单信息
func (r *ShopRequest) GetOrderInfo(id string) (model.Order, error) {
	var order model.Order
	err := r.DB.Where("id = ?", id).First(&order).Error
	return order, err
}
