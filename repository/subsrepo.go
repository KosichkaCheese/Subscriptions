package repository

import (
	"context"
	"subscriptions/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepoInterface interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetById(ctx context.Context, id uint) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id uint) error
	SumByFilters(ctx context.Context, userId, serviceName *string, start, end *time.Time) (int, error)
}

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) SubscriptionRepoInterface {
	return &SubscriptionRepo{db: db}
}

func (repo *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) error {
	return repo.db.WithContext(ctx).Create(subscription).Error
}

func (repo *SubscriptionRepo) GetById(ctx context.Context, id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := repo.db.WithContext(ctx).Preload("Service").First(&subscription, id).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (repo *SubscriptionRepo) GetAll(ctx context.Context) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := repo.db.WithContext(ctx).Preload("Service").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (repo *SubscriptionRepo) Update(ctx context.Context, subscription *models.Subscription) error {
	return repo.db.WithContext(ctx).Save(subscription).Error
}

func (repo *SubscriptionRepo) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&models.Subscription{}, id).Error
}

func (repo *SubscriptionRepo) SumByFilters(ctx context.Context, userId, serviceName *string, start, end *time.Time) (int, error) {
	query := repo.db.WithContext(ctx).Model(&models.Subscription{}).Select("SUM(Price)")

	if userId != nil {
		query = query.Where("user_id = ?", userId)
	}

	if serviceName != nil {
		query = query.Joins("JOIN services ON services.id = subscriptions.service_id").Where("services.name = ?", serviceName)
	}

	if start != nil {
		query = query.Where("start_date >= ?", start)
	}

	if end != nil {
		query = query.Where("end_date <= ? OR end_date IS NULL", end) //нужно учесть записи, у которых нет конца
	}

	var total int
	if err := query.Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
