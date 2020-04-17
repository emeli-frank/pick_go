package product

import "time"

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	RegularPrice float32 `json:"regular_price"`
	DiscountPrice float32 `json:"discount_price"`
}

type Order struct {
	Product Product `json:"product"`
	TimeOrdered time.Time `json:"time_ordered"`
}

type OrderHistory []Order

type Service interface {
	CreateProduct(product *Product) (int, error)
	GetProducts(index int, number int) ([]*Product, int, error)
	GetProduct(productId int, userId int) (product *Product, inCart bool, err error)
	GetCartItems(userId int) ([]*Product, error)
	GetOrderProducts(userId int) (*OrderHistory, error)
	SaveProductToCart(userId int, productId int) error
	SaveToOrderHistory(userId int, productIds []int, time time.Time) error
	DeleteProductFromCart(userId int, productIds []int) error
}

