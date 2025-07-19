package subscriptions_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"service/internal/subscriptions"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres password= dbname=sub sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	return db
}

func TestCreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := subscriptions.NewRepository(db)

	sub := subscriptions.Subscription{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
		EndDate:     nil,
	}

	err := repo.Create(context.Background(), sub)
	assert.NoError(t, err)

	// Получаем созданную подписку, предположим, что id автоинкремент
	subs, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, subs)

	got, err := repo.GetByID(context.Background(), subs[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, sub.ServiceName, got.ServiceName)
}
