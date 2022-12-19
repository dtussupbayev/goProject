package store

import (
	"context"
	"github.com/Assyl00/goProject/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Categories() CategoriesRepository
	Products() ProductsRepository
	Orders() OrdersRepository
	Reviews() ReviewsRepository
	Users() UsersRepository
}

type Authorization interface {
	CreateUser(user *models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type ProductsRepository interface {
	Create(ctx context.Context, product *models.Product) error
	All(ctx context.Context, filter *models.ProductsFilter) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type OrdersRepository interface {
	Create(ctx context.Context, order *models.Order) error
	All(ctx context.Context) ([]*models.Order, error)
	ByID(ctx context.Context, id int) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id int) error
}

type ReviewsRepository interface {
	Create(ctx context.Context, review *models.Review) error
	All(ctx context.Context, filter *models.ReviewsFilter) ([]*models.Review, error)
	ByID(ctx context.Context, id int) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id int) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}
