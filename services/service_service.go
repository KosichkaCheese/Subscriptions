package services

import (
	"context"
	"subscriptions/models"
	"subscriptions/repository"

	"go.uber.org/zap"
)

type ServiceServiceInterface interface {
	GetAll(ctx context.Context) ([]models.Service, error)
	Create(ctx context.Context, service *models.CreateService) (*models.Service, error)
	Delete(ctx context.Context, id uint) error
}

type ServiceService struct {
	repo   repository.ServiceRepoInterface
	logger *zap.SugaredLogger
}

func NewServiceService(repo repository.ServiceRepoInterface, logger *zap.SugaredLogger) ServiceServiceInterface {
	return &ServiceService{repo: repo, logger: logger}
}

func (s *ServiceService) GetAll(ctx context.Context) ([]models.Service, error) {
	res, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Errorf("GetAll services failed: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *ServiceService) Create(ctx context.Context, service *models.CreateService) (*models.Service, error) {
	newService := &models.Service{Name: service.Name}
	s.logger.Infof("Create service: %v", newService)
	err := s.repo.Create(ctx, newService)
	if err != nil {
		s.logger.Errorf("Create service failed: %v", err)
		return nil, err
	}
	return newService, nil
}

func (s *ServiceService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorf("Delete service failed: %v", err)
		return err
	}
	return nil
}
