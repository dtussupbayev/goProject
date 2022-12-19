package postgres

import (
	"github.com/Assyl00/goProject/internal/store"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	categories store.CategoriesRepository
	products   store.ProductsRepository
	orders     store.OrdersRepository
	reviews    store.ReviewsRepository
	users      store.UsersRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
