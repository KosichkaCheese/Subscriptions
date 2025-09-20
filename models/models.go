package models

import (
	"time"
)

type Service struct {
	ID        uint   `json:"id"`
	Name      string `gorm:"not null; unique" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Subscription struct {
	ID        uint       `json:"id"`
	ServiceID uint       `gorm:"not null; index" json:"service_id"`
	Service   Service    `gorm:"foreignkey:ServiceID" json:"service"`
	Price     uint       `gorm:"not null" json:"price"`
	UserID    string     `gorm:"type:uuid; not null; index" json:"user_id"`
	StartDate time.Time  `gorm:"not null" json:"start_date"`
	EndDate   *time.Time `json:"end_date"` //используем указатель, чтобы можно было использовать nil
	CreatedAt time.Time
	UpdatedAt time.Time
}

// модель для создания подписки
type CreateSubscription struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       *uint   `json:"price" binding:"required,gte=0"` //указатель чтобы отличать 0 от nil
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

// модель для обновления подписки
type UpdateSubscription struct {
	Price   *uint   `json:"price,omitempty"`
	EndDate *string `json:"end_date,omitempty"`
}

// модель для фильтрации
type SumFilter struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	StartDate   *string `form:"start_date"`
	EndDate     *string `form:"end_date"`
}

// модель для создания сервиса
type CreateService struct {
	Name string `json:"name" binding:"required"`
}
