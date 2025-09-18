package services

import (
	"context"
	"errors"
	"subscriptions/models"
	"subscriptions/repository"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidDate error

type SubscriptionServiceInterface interface {
	Create(ctx context.Context, subscription *models.CreateSubscription) (*models.Subscription, error)
	GetById(ctx context.Context, id uint) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	Update(ctx context.Context, id uint, subscription *models.UpdateSubscription) (*models.Subscription, error)
	Delete(ctx context.Context, id uint) error
	SumByFilters(ctx context.Context, filters *models.SumFilter) (int, error)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			service = &models.Service{Name: subscription.ServiceName}
			if err = s.servicerepo.Create(ctx, service); err != nil { //создание сервиса в базе данных
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
			ErrInvalidDate = errors.New("end date must be after start date")
			return nil, ErrInvalidDate
		}
	}

	sub := &models.Subscription{ServiceID: service.ID, UserID: subscription.UserID, StartDate: startDate, EndDate: endDate, Price: subscription.Price}
	err = s.subsrepo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetById(ctx context.Context, id uint) (*models.Subscription, error) {
	return s.subsrepo.GetById(ctx, id)
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]models.Subscription, error) {
	return s.subsrepo.GetAll(ctx)
}

func (s *SubscriptionService) Update(ctx context.Context, id uint, update *models.UpdateSubscription) (*models.Subscription, error) {
	sub, err := s.subsrepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if update.Price != nil {
		sub.Price = *update.Price
	}

	if update.EndDate != nil {
		endDate, err := time.Parse("01-2006", *update.EndDate)
		if err != nil {
			return nil, err
		}
		if endDate.Before(sub.StartDate) {
			ErrInvalidDate = errors.New("end date must be after start date")
			return nil, ErrInvalidDate
		}
		sub.EndDate = &endDate
	}

	err = s.subsrepo.Update(ctx, sub)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uint) error {
	return s.subsrepo.Delete(ctx, id)
}

func (s *SubscriptionService) SumByFilters(ctx context.Context, filters *models.SumFilter) (int, error) {
	var startDate, endDate *time.Time

	if filters == nil {
		return 0, errors.New("filters is nil")
	}

	if filters.StartDate != nil {
		start, err := time.Parse("01-2006", *filters.StartDate)
		if err != nil {
			return 0, err
		}
		startDate = &start
	}

	if filters.EndDate != nil {
		end, err := time.Parse("01-2006", *filters.EndDate)
		if err != nil {
			return 0, err
		}
		endDate = &end
	}

	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		ErrInvalidDate = errors.New("end date must be after start date")
		return 0, ErrInvalidDate
	}
	return s.subsrepo.SumByFilters(ctx, filters.UserID, filters.ServiceName, startDate, endDate)
}
