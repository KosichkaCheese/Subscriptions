package services

import (
	"context"
	"errors"
	"subscriptions/models"
	"subscriptions/repository"
	"time"

	"go.uber.org/zap"
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
	logger      *zap.SugaredLogger
}

func NewSubscriptionService(subsrepo repository.SubscriptionRepoInterface, servicerepo repository.ServiceRepoInterface, logger *zap.SugaredLogger) SubscriptionServiceInterface {
	return &SubscriptionService{subsrepo: subsrepo, servicerepo: servicerepo, logger: logger}
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *models.CreateSubscription) (*models.Subscription, error) {
	service, err := s.servicerepo.GetByName(ctx, subscription.ServiceName) //проверяем, есть ли сервис
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //если нет, добавляем
			service = &models.Service{Name: subscription.ServiceName}
			if err = s.servicerepo.Create(ctx, service); err != nil {
				s.logger.Errorf("Create service failed: %v", err)
				return nil, err
			}
		} else {
			s.logger.Errorf("GetByName service failed: %v", err)
			return nil, err
		}
	}

	startDate, err := time.Parse("01-2006", subscription.StartDate)
	if err != nil {
		s.logger.Errorf("Parsing start date failed: %v", err)
		return nil, err
	}

	var endDate *time.Time
	if subscription.EndDate != nil {
		endDateParse, err := time.Parse("01-2006", *subscription.EndDate)
		if err != nil {
			s.logger.Errorf("Parsing end date failed: %v", err)
			return nil, err
		}
		endDate = &endDateParse

		if endDate.Before(startDate) { //конец не должен быть раньше начала
			ErrInvalidDate = errors.New("end date must be after start date")
			s.logger.Error(ErrInvalidDate)
			return nil, ErrInvalidDate
		}
	}

	sub := &models.Subscription{ServiceID: service.ID, UserID: subscription.UserID, StartDate: startDate, EndDate: endDate, Price: *subscription.Price}
	s.logger.Infof("Creating subscription: %+v", sub)
	err = s.subsrepo.Create(ctx, sub)
	if err != nil {
		s.logger.Errorf("Create subscription failed: %v", err)
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetById(ctx context.Context, id uint) (*models.Subscription, error) {
	res, err := s.subsrepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorf("GetById subscription failed: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]models.Subscription, error) {
	res, err := s.subsrepo.GetAll(ctx)
	if err != nil {
		s.logger.Errorf("GetAll subscriptions failed: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *SubscriptionService) Update(ctx context.Context, id uint, update *models.UpdateSubscription) (*models.Subscription, error) {
	sub, err := s.subsrepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorf("GetById subscription failed: %v", err)
		return nil, err
	}

	if update.Price != nil {
		sub.Price = *update.Price
	}

	if update.EndDate != nil {
		endDate, err := time.Parse("01-2006", *update.EndDate)
		if err != nil {
			s.logger.Errorf("Parsing end date failed: %v", err)
			return nil, err
		}
		if endDate.Before(sub.StartDate) { //конец не должен быть раньше начала
			ErrInvalidDate = errors.New("end date must be after start date")
			s.logger.Error(ErrInvalidDate)
			return nil, ErrInvalidDate
		}
		sub.EndDate = &endDate
	}

	s.logger.Infof("Updating subscription: %+v", sub)

	err = s.subsrepo.Update(ctx, sub)
	if err != nil {
		s.logger.Errorf("Update subscription failed: %v", err)
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uint) error {
	err := s.subsrepo.Delete(ctx, id)
	if err != nil {
		s.logger.Errorf("Delete subscription failed: %v", err)
		return err
	}
	return nil
}

func (s *SubscriptionService) SumByFilters(ctx context.Context, filters *models.SumFilter) (int, error) {
	var startDate, endDate *time.Time

	if filters == nil {
		s.logger.Error("SumByFilters failed: filters is nil")
		return 0, errors.New("filters is nil")
	}

	if filters.StartDate != nil {
		start, err := time.Parse("01-2006", *filters.StartDate)
		if err != nil {
			s.logger.Errorf("Parsing start date failed: %v", err)
			return 0, err
		}
		startDate = &start
	}

	if filters.EndDate != nil {
		end, err := time.Parse("01-2006", *filters.EndDate)
		if err != nil {
			s.logger.Errorf("Parsing end date failed: %v", err)
			return 0, err
		}
		endDate = &end
	}

	if startDate != nil && endDate != nil && startDate.After(*endDate) { //конец не должен быть раньше начала
		ErrInvalidDate = errors.New("end date must be after start date")
		s.logger.Error(ErrInvalidDate)
		return 0, ErrInvalidDate
	}

	s.logger.Infof("SumByFilters: %+v", filters)

	res, err := s.subsrepo.SumByFilters(ctx, filters.UserID, filters.ServiceName, startDate, endDate)
	if err != nil {
		s.logger.Errorf("SumByFilters failed: %v", err)
		return 0, err
	}
	return res, nil
}
