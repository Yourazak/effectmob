package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, s Subscription) error
	GetByID(ctx context.Context, id int) (*Subscription, error)
	List(ctx context.Context) ([]Subscription, error)
	Update(ctx context.Context, s Subscription) error
	Delete(ctx context.Context, id int) error
	SumPrice(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

type Subscription struct {
	ID          int        `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}

func (r *PostgresRepository) Create(ctx context.Context, s Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*Subscription, error) {
	var s Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1`
	err := r.db.GetContext(ctx, &s, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("subscription not found")
	}
	return &s, err
}

func (r *PostgresRepository) List(ctx context.Context) ([]Subscription, error) {
	var subs []Subscription
	query := `SELECT * FROM subscriptions ORDER BY id`
	err := r.db.SelectContext(ctx, &subs, query)
	return subs, err
}

func (r *PostgresRepository) Update(ctx context.Context, s Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5
		WHERE id=$6
	`
	res, err := r.db.ExecContext(ctx, query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE id=$1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *PostgresRepository) SumPrice(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(price),0) FROM subscriptions
		WHERE user_id=$1 AND service_name=$2
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3)
	`
	err := r.db.GetContext(ctx, &total, query, userID, serviceName, startDate, endDate)
	return total, err
}
