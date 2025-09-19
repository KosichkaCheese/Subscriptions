package mocks

import (
	"context"
	"subscriptions/models"

	"github.com/stretchr/testify/mock"
)

type ServiceRepoMock struct {
	mock.Mock
}

func (s *ServiceRepoMock) Create(ctx context.Context, service *models.Service) error {
	args := s.Called(ctx, service)
	return args.Error(0)
}

func (s *ServiceRepoMock) GetByName(ctx context.Context, name string) (*models.Service, error) {
	args := s.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (s *ServiceRepoMock) GetAll(ctx context.Context) ([]models.Service, error) {
	args := s.Called(ctx)
	return args.Get(0).([]models.Service), args.Error(1)
}

func (s *ServiceRepoMock) GetById(ctx context.Context, id uint) (*models.Service, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*models.Service), args.Error(1)
}

func (s *ServiceRepoMock) Update(ctx context.Context, service *models.Service) error {
	args := s.Called(ctx, service)
	return args.Error(0)
}

func (s *ServiceRepoMock) Delete(ctx context.Context, id uint) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
