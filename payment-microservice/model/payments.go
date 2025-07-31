package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentSuccess   PaymentStatus = "SUCCESS"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	ID        uuid.UUID      `gorm:"type:char(255);primaryKey"`
	OrderID   uuid.UUID      `gorm:"not null"`
	UserID    uint           `gorm:"not null"`
	Amount    float64        `gorm:"not null"`
	Status    PaymentStatus  `gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
