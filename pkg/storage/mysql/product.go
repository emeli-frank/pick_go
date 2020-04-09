package mysql

import (
	"database/sql"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/domain/product"
	errors2 "github.com/emeli-frank/pick_go/pkg/errors"
)

type productStorage struct {
	DB *sql.DB
}

func NewProductStorage(db *sql.DB) *productStorage {
	return &productStorage{db}
}

func ListProducts() {

}

func (r *productStorage) SaveProduct(product *product.Product) (int, error) {
	const op = "productStorage.SaveProduct"

	fmt.Println(product)

	query := `INSERT INTO products (name, description, quantity, regular_price, discount_price) 
		VALUE (?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(query, product.Name, product.Description,
		product.Quantity, product.RegularPrice, product.DiscountPrice)
	if err != nil {
		return 0, errors2.Wrap(err, op, "inserting products")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors2.Wrap(err, op, "getting last insert id")
	}

	return int(id), nil
}

func (r *productStorage) GetProducts(index int, number int) ([]*product.Product, int, error) {
	const op = "productStorage.GetProducts"

	query := `SELECT id, name, description, regular_price, discount_price, quantity 
		FROM products
		LIMIT ?, ?`

	rows, err := r.DB.Query(query, index, number)
	if err != nil {
		return nil, 0, errors2.Wrap(err, op, "getting products")
	}
	defer rows.Close()

	var pp = []*product.Product{}

	for rows.Next() {
		p := &product.Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.RegularPrice, &p.DiscountPrice, &p.Quantity)
		if err := rows.Err(); err != nil {
			return nil, 0, errors2.Wrap(err, op, "scanning product into struct")
		}

		pp = append(pp, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors2.Wrap(err, op, "getting products")
	}

	// OBTAIN COUNT
	query = "SELECT COUNT(id) FROM products"
	row := r.DB.QueryRow(query)
	var total int
	err = row.Scan(&total)
	if err != nil {
		return nil, 0, errors2.Wrap(err, op, "getting count")
	}

	return pp, total, nil
}

func (r *productStorage) GetProduct(productId int, userId int) (*product.Product, bool, error) {
	const op = "productStorage.GetProduct"

	query := `SELECT id, name, description, regular_price, discount_price, quantity 
		FROM products WHERE id = ?`
	row := r.DB.QueryRow(query, productId)
	p := &product.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.RegularPrice, &p.DiscountPrice, &p.Quantity)
	if err == sql.ErrNoRows {
		return nil, false, errors2.Wrap(&errors2.NotFound{Err:err}, op, "getting and scanning rows into variable")
	}
	if err != nil {
		return nil, false, errors2.Wrap(err, op, "getting and scanning rows into variable")
	}

	var inCart bool = true
	query = `SELECT product_id FROM cart_items WHERE product_id = ? AND user_id = ?`
	row = r.DB.QueryRow(query, productId, userId)
	var discard interface{}
	err = row.Scan(&discard)
	if err == sql.ErrNoRows {
		inCart = false
	} else if err != nil {
		return nil, false, errors2.Wrap(err, op, "getting and scanning rows into variable")
	}

	return p, inCart, nil
}

func (r *productStorage) GetCartItems(userId int) ([]*product.Product, error) {
	const op = "productStorage.GetCartItems"

	query := `SELECT products.id, name, description, regular_price, discount_price, quantity
		FROM products
		LEFT JOIN cart_items
			ON products.id = cart_items.product_id
		WHERE cart_items.user_id = ?`

	rows, err := r.DB.Query(query, userId)
	if err != nil {
		return nil, errors2.Wrap(err, op, "running query")
	}
	defer rows.Close()

	var pp = []*product.Product{}

	for rows.Next() {
		p := &product.Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.RegularPrice, &p.DiscountPrice, &p.Quantity)
		if err := rows.Err(); err != nil {
			return nil, errors2.Wrap(err, op, "scanning product into struct")
		}

		pp = append(pp, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errors2.Wrap(err, op, "getting products")
	}

	fmt.Println(pp)

	return pp, nil
}

func (r *productStorage) GetOrderProducts(userId int) ([]*product.Product, error) {
	const op = "productStorage.GetCartItems"

	query := `SELECT products.id, name, description, regular_price, discount_price, quantity
		FROM products
		LEFT JOIN order_history
			ON products.id = order_history.product_id
		WHERE order_history.user_id = ?`

	rows, err := r.DB.Query(query, userId)
	if err != nil {
		return nil, errors2.Wrap(err, op, "running query")
	}
	defer rows.Close()

	var pp = []*product.Product{}

	for rows.Next() {
		p := &product.Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.RegularPrice, &p.DiscountPrice, &p.Quantity)
		if err := rows.Err(); err != nil {
			return nil, errors2.Wrap(err, op, "scanning product into struct")
		}

		pp = append(pp, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errors2.Wrap(err, op, "getting products")
	}

	return pp, nil
}
