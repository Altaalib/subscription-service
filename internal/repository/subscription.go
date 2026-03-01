package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetAll(ctx context.Context) ([]*model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input model.UpdateSubscriptionInput) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string) (*model.TotalCostResponse, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error) {
	sub := &model.Subscription{}
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query,
		input.ServiceName,
		input.Price,
		input.UserID,
		input.StartDate,
		input.EndDate,
	).StructScan(sub)
	if err != nil {
		return nil, fmt.Errorf("repository.Create: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	sub := &model.Subscription{}
	query := `SELECT * FROM subscriptions WHERE id = $1`

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(sub)
	if err != nil {
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepo) GetAll(ctx context.Context) ([]*model.Subscription, error) {
	var subs []*model.Subscription
	query := `SELECT * FROM subscriptions ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &subs, query)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	return subs, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, input model.UpdateSubscriptionInput) (*model.Subscription, error) {
	sub := &model.Subscription{}
	query := `
		UPDATE subscriptions
		SET
			service_name = COALESCE($1, service_name),
			price        = COALESCE($2, price),
			start_date   = COALESCE($3, start_date),
			end_date     = COALESCE($4, end_date),
			updated_at   = NOW()
		WHERE id = $5
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query,
		input.ServiceName,
		input.Price,
		input.StartDate,
		input.EndDate,
		id,
	).StructScan(sub)
	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}
	return sub, nil
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}
	return nil
}

func (r *subscriptionRepo) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string) (*model.TotalCostResponse, error) {
	result := &model.TotalCostResponse{}
	query := `
		SELECT COALESCE(SUM(price), 0) as total_cost
		FROM subscriptions
		WHERE user_id = $1 AND service_name = $2`

	err := r.db.QueryRowxContext(ctx, query, userID, serviceName).Scan(&result.TotalCost)
	if err != nil {
		return nil, fmt.Errorf("repository.GetTotalCost: %w", err)
	}
	return result, nil
}
