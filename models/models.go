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
	EndDate   *time.Time `json:"end_date"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateSubscription struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       uint    `json:"price" binding:"required,gte=0"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}
