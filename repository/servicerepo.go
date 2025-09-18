package repository

import (
	"context"
	"subscriptions/models"

	"gorm.io/gorm"
)

type ServiceRepoInterface interface {
	Create(ctx context.Context, service *models.Service) error
	GetAll(ctx context.Context) ([]models.Service, error)
	GetById(ctx context.Context, id uint) (*models.Service, error)
	GetByName(ctx context.Context, name string) (*models.Service, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id uint) error
}

type ServiceRepo struct {
	db *gorm.DB
}

func NewServiceRepo(db *gorm.DB) ServiceRepoInterface { //создание репозитория для сервисов
	return &ServiceRepo{db: db}
}

func (repo *ServiceRepo) Create(ctx context.Context, service *models.Service) error { //создание сервиса
	return repo.db.WithContext(ctx).Create(service).Error
}

func (repo *ServiceRepo) GetAll(ctx context.Context) ([]models.Service, error) { //получение всех сервисов
	var services []models.Service
	if err := repo.db.WithContext(ctx).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (repo *ServiceRepo) GetById(ctx context.Context, id uint) (*models.Service, error) { //получение сервиса по id
	var service models.Service
	if err := repo.db.WithContext(ctx).First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (repo *ServiceRepo) GetByName(ctx context.Context, name string) (*models.Service, error) { //получение сервиса по названию
	var service models.Service
	if err := repo.db.WithContext(ctx).Where("name = ?", name).First(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (repo *ServiceRepo) Update(ctx context.Context, service *models.Service) error { //обновление сервиса
	return repo.db.WithContext(ctx).Save(service).Error
}

func (repo *ServiceRepo) Delete(ctx context.Context, id uint) error { //удаление сервиса
	return repo.db.WithContext(ctx).Delete(&models.Service{}, id).Error
}
