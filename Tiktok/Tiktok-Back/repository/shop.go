package repository

import (
	"Tiktok/model"
	"Tiktok/repository/list"
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

// 获取订单
func (r *ShopRequest) GetOrderInfo(OrderId int64) (model.Order, error) {
	var Order model.Order
	Order.ID = OrderId
	result := r.DB.First(&Order)
	if result.Error != nil {
		return Order, result.Error
	}
	return Order, nil

}

// 创建订单
func (r *ShopRequest) CreateOrder(Order model.Order) error {
	return r.DB.Create(Order).Error
}

// 搜索订单
func (r *ShopRequest) SearchOrder(KeyWord, OrderBy string) ([]model.Order, error) {
	var Orders []model.Order
	Orders, _, err := list.Query[model.Order](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("is_delete = ?", false).Select("%" + KeyWord + "%")
		},
		Order: OrderBy,

		//Preloads: []string{"Author"},
		PageInfo: list.PageInfo{Page: 1, Limit: 10},
	})
	return Orders, err
}

// 更新订单
func (r *ShopRequest) UpdateOrder(Order model.Order) error {
	return r.DB.Updates(Order).Error
}

// 获取商品
func (r *ShopRequest) GetGoodInfo(GoodId int64) (model.Good, error) {
	var Good model.Good
	Good.ID = GoodId
	result := r.DB.First(&Good)
	if result.Error != nil {
		return Good, result.Error
	}
	return Good, nil

}

// 搜索商品
func (r *ShopRequest) SearchGood(GoodsCategoryId int64, KeyWord, OrderBy string) ([]model.Good, error) {

	var Goods []model.Good
	Goods, _, err := list.Query[model.Good](list.Options{
		Where: func(db *gorm.DB) *gorm.DB {
			return db.Where("category_id = ? AND is_delete = ?", GoodsCategoryId, false).Select("%" + KeyWord + "%")
		},
		Order: OrderBy,

		//Preloads: []string{"Author"},
		PageInfo: list.PageInfo{Page: 1, Limit: 10},
	})
	return Goods, err

}

// 创建商品
func (r *ShopRequest) CreateGood(Good model.Good) error {
	return r.DB.Create(Good).Error
}

// 更新商品
func (r *ShopRequest) UpdateGood(Good model.Good) error {
	return r.DB.Updates(Good).Error
}

// 获取购物车
func (r *ShopRequest) GetCartInfo(CartId int64) (model.Cart, error) {
	var Cart model.Cart
	Cart.ID = CartId
	result := r.DB.First(&Cart)
	if result.Error != nil {
		return Cart, result.Error
	}
	return Cart, nil

}

// 创建购物车
func (r *ShopRequest) CreateCart(Cart model.Cart) error {
	return r.DB.Create(Cart).Error
}

//
