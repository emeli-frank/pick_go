package product

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Quantity int `json:"quantity"`
	RegularPrice float32 `json:"regular_price"`
	DiscountPrice float32 `json:"discount_price"`
}

type Service interface {
	CreateProduct(product *Product) (int, error)
	GetProducts(index int, number int) ([]*Product, int, error)
	GetProduct(productId int, userId int) (product *Product, inCart bool, err error)
	GetCartItems(userId int) ([]*Product, error)
	GetOrderProducts(userId int) ([]*Product, error)
}

