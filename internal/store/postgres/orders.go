package postgres

import (
	"context"
	"github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Orders() store.OrdersRepository {
	if db.orders == nil {
		db.orders = NewOrdersRepository(db.conn)
	}

	return db.orders
}

type OrdersRepository struct {
	conn *sqlx.DB
}

func NewOrdersRepository(conn *sqlx.DB) store.OrdersRepository {
	return &OrdersRepository{conn: conn}
}

func (o *OrdersRepository) Create(ctx context.Context, order *models.Order) error {
	_, err := o.conn.ExecContext(ctx, "INSERT INTO orders(user_id) VALUES ($1)", order.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrdersRepository) All(ctx context.Context) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	if err := o.conn.Select(&orders, "SELECT * FROM orders"); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrdersRepository) AllOfUsers(ctx context.Context, userId int) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	err := o.conn.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrdersRepository) ByID(ctx context.Context, id int) (*models.Order, error) {
	order := new(models.Order)

	if err := o.conn.Get(order, "SELECT id, user_id FROM categories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrdersRepository) Update(ctx context.Context, order *models.Order) error {
	//TODO implement me
	panic("implement me")
}

func (o *OrdersRepository) Delete(ctx context.Context, id int) error {
	_, err := o.conn.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
