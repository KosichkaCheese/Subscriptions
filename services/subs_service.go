package services

import (
	"context"
	"errors"
	"subscriptions/models"
	"subscriptions/repository"
	"time"

	"gorm.io/gorm"
)

type SubscriptionServiceInterface interface {
	Create(ctx context.Context, subscription *models.CreateSubscription) (*models.Subscription, error)
	GetById(ctx context.Context, id uint) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id uint) error
	SumByFilters(ctx context.Context, userId, serviceName string, start, end *string) (int, error)
}

type SubscriptionService struct {
	subsrepo    repository.SubscriptionRepoInterface
	servicerepo repository.ServiceRepoInterface
}

func NewSubscriptionService(subsrepo repository.SubscriptionRepoInterface, servicerepo repository.ServiceRepoInterface) SubscriptionServiceInterface {
	return &SubscriptionService{subsrepo: subsrepo, servicerepo: servicerepo}
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *models.CreateSubscription) (*models.Subscription, error) {
	service, err := s.servicerepo.GetByName(ctx, subscription.ServiceName)
	if err != nil { //проверка на наличие сервиса в базе данных
		if err == gorm.ErrRecordNotFound {
			if err = s.servicerepo.Create(ctx, &models.Service{Name: subscription.ServiceName}); err != nil { //создание сервиса в базе данных
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	startDate, err := time.Parse("01-2006", subscription.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if subscription.EndDate != nil {
		endDateParse, err := time.Parse("01-2006", *subscription.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &endDateParse

		if endDate.Before(startDate) {
			ErrInvalidDate := errors.New("end date must be after start date")
			return nil, ErrInvalidDate
		}
	}

	newSubscription := &models.Subscription{ServiceID: service.ID, UserID: subscription.UserID, StartDate: startDate, EndDate: endDate, Price: subscription.Price}
	err = s.subsrepo.Create(ctx, newSubscription)
	if err != nil {
		return nil, err
	}
	return newSubscription, nil
}

func (s *SubscriptionService) GetById(ctx context.Context, id uint) (*models.Subscription, error) {
	return s.subsrepo.GetById(ctx, id)
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]models.Subscription, error) {
	return s.subsrepo.GetAll(ctx)
}

func (s *SubscriptionService) Update(ctx context.Context, subscription *models.Subscription) error {
	return s.subsrepo.Update(ctx, subscription)
}

func (s *SubscriptionService) Delete(ctx context.Context, id uint) error {
	return s.subsrepo.Delete(ctx, id)
}

func (s *SubscriptionService) SumByFilters(ctx context.Context, userId, serviceName string, start, end *string) (int, error) {
	return s.subsrepo.SumByFilters(ctx, userId, serviceName, start, end)
}
