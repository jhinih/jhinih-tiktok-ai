package types

type Order struct {
	ID          string `json:"id"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	OwnerID     string `json:"owner_id"`
	IsDelete    bool   `json:"is_delete"`
	Price       string `json:"price"`
	GoodID      string `json:"good_id"`
	AddressID   string `json:"address_id"`
	Status      string `json:"status"`
	Pay         bool   `json:"pay"`
	PayBy       string `json:"pay_by"`
	PayTime     string `json:"pay_time"`
}

// 获取订单
type GetOrderInfoRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	OwnerID string `json:"owner_id" binding:"-"`
}
type GetOrderInfoResponse struct {
	Order Order `json:"order"`
}

// 创建订单
type CreateOrderRequest struct {
	OwnerID   string `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	Price     string `json:"price" gorm:"column:price;type:bigint;type:bigint"`
	AddressID string `json:"address_id"`
}
type CreateOrderResponse struct {
}

// 创建订单
type UpdateOrderRequest struct {
	Price     string `json:"price"`
	AddressID string `json:"address_id"`
	GoodID    string `json:"good_id"`
}
type UpdateOrderResponse struct {
}

// 创建订单
type SearchOrderRequest struct {
	KeyWord string `json:"key_word" binding:"-"`
	OrderBy string `json:"order_by" binding:"-"`
}
type SearchOrderResponse struct {
	Orders []Order `json:"orders"`
}

type Good struct {
	ID          string `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	OrderID     string `json:"order_id" gorm:"column:order_id;type:bigint"`
	CartID      string `json:"cart_id" gorm:"column:cart_id;type:bigint"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	Name        string `json:"name" gorm:"column:name"`
	Info        string `json:"info" gorm:"column:info"`
	Category    string `json:"category" gorm:"column:category"`
	Tags        string `json:"tags" gorm:"column:tags"`
	Img         string `json:"img" gorm:"column:img"`

	OriginPrice string `json:"originPrice" gorm:"column:origin_price"`
	SellPrice   string `json:"sellPrice" gorm:"column:sell_price"`

	Stock  string `json:"stock" gorm:"column:stock"`
	Status string `json:"status" gorm:"column:status"`

	OwnerID string `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
}

type GetGoodInfoRequest struct {
	GoodID  string `json:"good_id" binding:"required"`
	OwnerID string `json:"owner_id" binding:"-"`
}

// 获取商品响应
type GetGoodInfoResponse struct {
	Good Good `json:"good"`
}

// 创建商品
type CreateGoodRequest struct {
	OwnerID  string `json:"owner_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Info     string `json:"info" binding:"required"`
	Category string `json:"category" binding:"required"`
	Tags     string `json:"tags" binding:"required"`
	Img      string `json:"img" binding:"required"`

	OriginPrice string `json:"origin_price" binding:"required"`
	SellPrice   string `json:"sell_price" binding:"required"`
	Stock       string `json:"stock" binding:"required"`
}
type CreateGoodResponse struct {
}

// 搜索商品
type SearchGoodRequest struct {
	GoodsCategoryId string `json:"goods_category_id" binding:"required"`
	KeyWord         string `json:"key_word" binding:"-"`
	OrderBy         string `json:"order_by" binding:"-"`
}
type SearchGoodResponse struct {
	Goods []Good `json:"goods"`
}

// 更新商品
type UpdateGoodRequest struct {
	Name       string `json:"name" gorm:"column:name"`
	Info       string `json:"info" gorm:"column:info"`
	Category   string `json:"category" gorm:"column:category"`
	CategoryId string `json:"category_id"`
	Tags       string `json:"tags" gorm:"column:tags"`
	Img        string `json:"img" gorm:"column:img"`

	OriginPrice string `json:"originPrice" gorm:"column:origin_price"`
	SellPrice   string `json:"sellPrice" gorm:"column:sell_price"`

	Stock  string `json:"stock" gorm:"column:stock"`
	Status string `json:"status" gorm:"column:status"`
}
type UpdateGoodResponse struct {
}

type Cart struct {
	ID          string `json:"id"`
	CreatedTime string `gorm:"column:created_time;type:bigint"`
	UpdatedTime string `gorm:"column:updated_time;type:bigint"`
	OwnerID     string `json:"owner_id"`
	GoodsID     string `json:"goods_id"`
	GoodsNum    string `json:"goods_num"`
	IsDelete    bool   `json:"is_delete"`
}
type GetCartInfoRequest struct {
	CartID  string `json:"cart_id" binding:"required"`
	OwnerID string `json:"owner_id" binding:"-"`
}

// 获取购物车响应
type GetCartInfoResponse struct {
	Cart Cart `json:"cart"`
}

type CreateCartRequest struct {
	OwnerID string `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	Goods   []Good `json:"goods" binding:"required"`
}

// 创建购物车响应
type CreateCartResponse struct {
}
