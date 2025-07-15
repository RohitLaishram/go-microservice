package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Paid      OrderStatus = "PAID"
	Shipped   OrderStatus = "SHIPPED"
	Delivered OrderStatus = "DELIVERED"
	Canceled  OrderStatus = "CANCELED"
)

type Order struct {
	ID          uuid.UUID      `gorm:"type:char(255);primaryKey"`
	UserID      uint           `gorm:"not null"`
	ProductID   string         `gorm:"not null"`
	Quantity    int            `gorm:"not null;default:1"`
	TotalAmount float64        `gorm:"not null"` // price * quantity
	Status      OrderStatus    `gorm:"type:varchar(20);default:'PENDING'"`
	PaymentID   *uint          `gorm:"default:null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
