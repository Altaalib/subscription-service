package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type SubscriptionService interface {
	Create(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetAll(ctx context.Context) ([]*model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input model.UpdateSubscriptionInput) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string) (*model.TotalCostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error) {
	sub, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}
	return sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return sub, nil
}

func (s *subscriptionService) GetAll(ctx context.Context) ([]*model.Subscription, error) {
	subs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetAll: %w", err)
	}
	return subs, nil
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, input model.UpdateSubscriptionInput) (*model.Subscription, error) {
	sub, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}
	return sub, nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}
	return nil
}

func (s *subscriptionService) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string) (*model.TotalCostResponse, error) {
	result, err := s.repo.GetTotalCost(ctx, userID, serviceName)
	if err != nil {
		return nil, fmt.Errorf("service.GetTotalCost: %w", err)
	}
	return result, nil
}
