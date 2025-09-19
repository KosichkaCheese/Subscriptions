package mocks

import (
	"context"
	"subscriptions/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type SubscriptionRepoMock struct {
	mock.Mock
}

func (s *SubscriptionRepoMock) Create(ctx context.Context, subscription *models.Subscription) error {
	args := s.Called(ctx, subscription)
	return args.Error(0)
}

func (s *SubscriptionRepoMock) GetById(ctx context.Context, id uint) (*models.Subscription, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (s *SubscriptionRepoMock) GetAll(ctx context.Context) ([]models.Subscription, error) {
	args := s.Called(ctx)
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (s *SubscriptionRepoMock) Update(ctx context.Context, subscription *models.Subscription) error {
	args := s.Called(ctx, subscription)
	return args.Error(0)
}

func (s *SubscriptionRepoMock) Delete(ctx context.Context, id uint) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func (s *SubscriptionRepoMock) SumByFilters(ctx context.Context, userId, serviceName *string, start, end *time.Time) (int, error) {
	args := s.Called(ctx, userId, serviceName, start, end)
	return args.Get(0).(int), args.Error(1)
}
