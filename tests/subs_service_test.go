package tests

import (
	"context"
	"subscriptions/services"
	"testing"
	"time"

	"subscriptions/models"
	"subscriptions/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestCreate_NoService(t *testing.T) { //нет сервиса с таким именем
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	price := uint(500)
	createSub := &models.CreateSubscription{
		ServiceName: "Spotify",
		UserID:      "6a2995b1-9967-473c-ab26-2710f6e66fd5",
		Price:       &price,
		StartDate:   "01-2025",
	}

	srepo.On("GetByName", ctx, "Spotify").Return(nil, gorm.ErrRecordNotFound)
	srepo.On("Create", ctx, mock.AnythingOfType("*models.Service")).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Service)
		s.ID = 1 // отдаем сервису id 1, чтобы функция не упала
	}).Return(nil)
	subrepo.On("Create", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil)

	res, err := subService.Create(ctx, createSub)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotZero(t, res.ServiceID)
	assert.Equal(t, "6a2995b1-9967-473c-ab26-2710f6e66fd5", res.UserID)
	assert.Equal(t, uint(500), res.Price)
	assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), res.StartDate)

	srepo.AssertExpectations(t)
	subrepo.AssertExpectations(t)
}

func TestCreate_InvalidDate(t *testing.T) { //невалидная дата
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	end := "01-2024"
	price := uint(500)
	createSub := &models.CreateSubscription{
		ServiceName: "Spotify",
		UserID:      "6a2995b1-9967-473c-ab26-2710f6e66fd5",
		Price:       &price,
		StartDate:   "01-2025",
		EndDate:     &end,
	}

	srepo.On("GetByName", ctx, "Spotify").Return(&models.Service{ID: 1, Name: "Spotify"}, nil)
	subrepo.On("Create", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil)

	res, err := subService.Create(ctx, createSub)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "end date must be after start date")
}

func TestUpdate_Success(t *testing.T) { //успешное обновление
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	existedSub := &models.Subscription{
		ID:        1,
		ServiceID: 1,
		UserID:    "6a2995b1-9967-473c-ab26-2710f6e66fd5",
		Price:     500,
		StartDate: time.Now(),
		EndDate:   nil,
	}

	newPrice := uint(600)

	subrepo.On("GetById", ctx, uint(1)).Return(existedSub, nil)
	subrepo.On("Update", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil)

	res, err := subService.Update(ctx, 1, &models.UpdateSubscription{Price: &newPrice})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(600), res.Price)
}

func TestUpdate_InvalidDate(t *testing.T) { //невалидная дата
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	existedSub := &models.Subscription{
		ID:        1,
		ServiceID: 1,
		UserID:    "6a2995b1-9967-473c-ab26-2710f6e66fd5",
		Price:     500,
		StartDate: time.Now(),
		EndDate:   nil,
	}

	newEnd := "01-2024"

	subrepo.On("GetById", ctx, uint(1)).Return(existedSub, nil)
	subrepo.On("Update", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil)

	res, err := subService.Update(ctx, 1, &models.UpdateSubscription{EndDate: &newEnd})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "end date must be after start date")
}

func TestSumByFilters_Success(t *testing.T) { //успешное получение суммы
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	userID := "6a2995b1-9967-473c-ab26-2710f6e66fd5"
	serviceName := "Spotify"
	start := "01-2025"
	end := "01-2026"
	filters := &models.SumFilter{
		UserID:      &userID,
		ServiceName: &serviceName,
		StartDate:   &start,
		EndDate:     &end,
	}

	endDate, _ := time.Parse("01-2006", end)
	startDate, _ := time.Parse("01-2006", start)

	subrepo.On("SumByFilters", ctx, &userID, &serviceName, &startDate, &endDate).Return(1000, nil)

	res, err := subService.SumByFilters(ctx, filters)
	assert.NoError(t, err)
	assert.NotZero(t, res)
	assert.Equal(t, 1000, res)

}

func TestSumByFilters_InvalidDate(t *testing.T) { //невалидная дата
	ctx := context.Background()
	srepo := new(mocks.ServiceRepoMock)
	subrepo := new(mocks.SubscriptionRepoMock)
	log := zap.NewNop().Sugar()

	subService := services.NewSubscriptionService(subrepo, srepo, log)

	userID := "6a2995b1-9967-473c-ab26-2710f6e66fd5"
	serviceName := "Spotify"
	start := "01-2025"
	end := "01-2024"
	filters := &models.SumFilter{
		UserID:      &userID,
		ServiceName: &serviceName,
		StartDate:   &start,
		EndDate:     &end,
	}

	endDate, _ := time.Parse("01-2006", end)
	startDate, _ := time.Parse("01-2006", start)

	subrepo.On("SumByFilters", ctx, &userID, &serviceName, &startDate, &endDate).Return(1000, nil)

	res, err := subService.SumByFilters(ctx, filters)
	assert.Zero(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "end date must be after start date")
}
