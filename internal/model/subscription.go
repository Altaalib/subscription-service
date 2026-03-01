package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ServiceName string    `db:"service_name" json:"service_name" validate:"required"`
	Price       int       `db:"price" json:"price" validate:"required,min=1"`
	UserID      uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	StartDate   string    `db:"start_date" json:"start_date" validate:"required"`
	EndDate     *string   `db:"end_date" json:"end_date,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateSubscriptionInput struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
