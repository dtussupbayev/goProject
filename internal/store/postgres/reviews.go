package postgres

import (
	"context"
	"fmt"
	"github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db DB) Reviews() store.ReviewsRepository {
	if db.reviews == nil {
		db.reviews = NewReviewsRepository(db.conn)
	}

	return db.reviews
}

type ReviewsRepository struct {
	conn *sqlx.DB
}

func NewReviewsRepository(conn *sqlx.DB) store.ReviewsRepository {
	return &ReviewsRepository{conn: conn}
}

func (p *ReviewsRepository) Create(ctx context.Context, review *models.Review) error {
	_, err := p.conn.ExecContext(ctx, "INSERT INTO reviews(stars, body, product_id, product_name) VALUES ($1, $2, $3, $4)",
		review.Stars, review.Body, review.ProductID, review.ProductName)
	if err != nil {
		return err
	}

	return nil
}

func (p *ReviewsRepository) All(ctx context.Context, filter *models.ReviewsFilter) ([]*models.Review, error) {
	reviews := make([]*models.Review, 0)

	basicQuery := "SELECT * FROM reviews"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := p.conn.Select(&reviews, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return reviews, nil
	}

	if err := p.conn.Select(&reviews, basicQuery); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (p *ReviewsRepository) ByID(ctx context.Context, id int) (*models.Review, error) {
	reviews := new(models.Review)

	if err := p.conn.Get(reviews, "SELECT * FROM reviews WHERE id=$1", id); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (p *ReviewsRepository) ByProductID(ctx context.Context, productID int) ([]*models.Review, error) {
	reviews := make([]*models.Review, 0)

	err := p.conn.SelectContext(ctx, &reviews, "SELECT * FROM reviews WHERE product_id=$1", productID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (p *ReviewsRepository) Update(ctx context.Context, review *models.Review) error {
	_, err := p.conn.Exec("UPDATE reviews SET stars = $1, body = $2 WHERE id = $3", review.Stars, review.Body, review.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *ReviewsRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM reviews WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
