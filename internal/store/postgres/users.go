package postgres

import (
	"context"
	"github.com/dtusupbaev/goProject/internal/models"
	"github.com/dtusupbaev/goProject/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Users() store.UsersRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}

	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (u UsersRepository) Create(ctx context.Context, user *models.User) error {
	//if err := user.Validate(); err != nil {
	//	return err
	//}
	//
	//if err := user.BeforeCreating(); err != nil {
	//	return err
	//}

	_, err := u.conn.ExecContext(ctx, "INSERT INTO users(email, password) VALUES ($1, $2)",
		user.Email, user.EncryptedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) All(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)

	if err := u.conn.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UsersRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	//if err := user.Validate(); err != nil {
	//	return err
	//}
	//
	//if err := user.BeforeCreating(); err != nil {
	//	return err
	//}

	_, err := u.conn.ExecContext(ctx, "UPDATE users SET password = $2 WHERE id=$8",
		user.EncryptedPassword, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id int) error {
	if _, err := u.conn.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
