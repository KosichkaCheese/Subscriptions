package services

import (
	"context"
	"subscriptions/models"
	"subscriptions/repository"
)

type ServiceServiceInterface interface {
	GetAll(ctx context.Context) ([]models.Service, error)
	Create(ctx context.Context, service *models.CreateService) (*models.Service, error)
}

type ServiceService struct {
	repo repository.ServiceRepoInterface
}

func NewServiceService(repo repository.ServiceRepoInterface) ServiceServiceInterface {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll(ctx context.Context) ([]models.Service, error) {
	return s.repo.GetAll(ctx)
}

func (s *ServiceService) Create(ctx context.Context, service *models.CreateService) (*models.Service, error) {
	newService := &models.Service{Name: service.Name}
	err := s.repo.Create(ctx, newService)
	if err != nil {
		return nil, err
	}
	return newService, nil
}
