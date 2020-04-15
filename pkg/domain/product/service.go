package product

import (
	"github.com/emeli-frank/pick_go/pkg/errors"
	"time"
)

type repository interface {
	SaveProduct(product *Product) (int, error)
	GetProducts(index int, number int) ([]*Product, int, error)
	GetProduct(productId int, userId int) (product *Product, inCart bool, err error)
	GetCartItems(userId int) ([]*Product, error)
	GetOrderProducts(userId int) ([]*Product, error)

	SaveProductToCart(userId int, productIds []int) error
	SaveToOrderHistory(userId int, productIds []int, time time.Time) error

	DeleteProductFromCart(userId int, productIds []int) error
}

func New(repo repository) *service {
	return &service{
		r: repo,
	}
}

type service struct {
	r repository
}

func (s *service) CreateProduct(product *Product) (int, error) {
	const op = "productService.CreateProduct"

	id, err := s.r.SaveProduct(product)
	if err != nil {
		return id, errors.Wrap(err, op, "calling repo to save products")
	}

	return id, nil
}

func (s *service) GetProducts(index int, number int) ([]*Product, int, error) {
	const op = "productService.GetProducts"

	pp, total, err := s.r.GetProducts(index, number)
	if err != nil {
		return nil, 0, errors.Wrap(err, op, "getting products from repo")
	}

	return pp, total, nil
}

func (s *service) GetProduct(productId int, userId int) (*Product, bool, error) {
	const op = "productService.GetProduct"

	p, inCart, err := s.r.GetProduct(productId, userId)
	if err != nil {
		return nil, false, errors.Wrap(err, op, "getting product from repo")
	}

	return p, inCart, nil
}

func (s *service) GetCartItems(userId int) ([]*Product, error) {
	const op = "productService.GetCartItems"

	pp, err := s.r.GetCartItems(userId)
	if err != nil {
		return nil, errors.Wrap(err, op, "getting products from repo")
	}

	return pp, nil
}

func (s *service) SaveProductToCart(userId int, productIds []int) error {
	const op = "productService.SaveProductToCart"

	err := s.r.SaveProductToCart(userId, productIds)
	if err != nil {
		return errors.Wrap(err, op, "calling repo to save products")
	}

	return nil
}

func (s *service) GetOrderProducts(userId int) ([]*Product, error) {
	const op = "productService.GetOrderProducts"

	pp, err := s.r.GetOrderProducts(userId)
	if err != nil {
		return nil, errors.Wrap(err, op, "getting products from repo")
	}

	return pp, nil
}

func (s *service) SaveToOrderHistory(userId int, productIds []int, time time.Time) error {
	const op = "productService.SaveProductToCart"

	err := s.r.SaveToOrderHistory(userId, productIds, time)
	if err != nil {
		return errors.Wrap(err, op, "calling repo to save products")
	}

	/* error ignored for now */
	_ = s.DeleteProductFromCart(userId, productIds)

	return nil
}

func (s *service) DeleteProductFromCart(userId int, productIds []int) error {
	const op = "DeleteProductFromCart.SaveProductToCart"

	err := s.r.DeleteProductFromCart(userId, productIds)
	if err != nil {
		return errors.Wrap(err, op, "deleting product from repo")
	}

	return nil
}
