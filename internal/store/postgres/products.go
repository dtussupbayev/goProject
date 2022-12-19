package postgres

import (
	"context"
	"fmt"
	"github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db DB) Products() store.ProductsRepository {
	if db.products == nil {
		db.products = NewProductsRepository(db.conn)
	}

	return db.products
}

type ProductsRepository struct {
	conn *sqlx.DB
}

func NewProductsRepository(conn *sqlx.DB) store.ProductsRepository {
	return &ProductsRepository{conn: conn}
}

func (p *ProductsRepository) Create(ctx context.Context, product *models.Product) error {
	_, err := p.conn.ExecContext(ctx, "INSERT INTO products(name, description, price, category_id) VALUES ($1, $2, $3, $4)",
		product.Name, product.Description, product.Price, product.CategoryID)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsRepository) All(ctx context.Context, filter *models.ProductsFilter) ([]*models.Product, error) {
	products := make([]*models.Product, 0)

	basicQuery := "SELECT * FROM products"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := p.conn.Select(&products, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return products, nil
	}

	if err := p.conn.Select(&products, basicQuery); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductsRepository) ByID(ctx context.Context, id int) (*models.Product, error) {
	products := new(models.Product)

	if err := p.conn.Get(products, "SELECT * FROM products WHERE id=$1", id); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductsRepository) ByCategoryID(ctx context.Context, categoryID int) ([]*models.Product, error) {
	products := make([]*models.Product, 0)

	err := p.conn.SelectContext(ctx, &products, "SELECT * FROM products WHERE category_id=$1", categoryID)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductsRepository) Update(ctx context.Context, product *models.Product) error {
	_, err := p.conn.Exec("UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4", product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductsRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
