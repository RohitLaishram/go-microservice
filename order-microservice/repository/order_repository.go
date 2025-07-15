package repository

import (
	"order/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order model.Order) (*model.Order, error) {
	order.ID = uuid.New()

	err := r.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Get(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := r.DB.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.DB.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}
