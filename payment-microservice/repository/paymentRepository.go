package repository

import (
	"payment/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}
func (r *PaymentRepository) CreatePayment(payment *model.Payment) error {
	payment.ID = uuid.New()
	return r.db.Create(payment).Error
}
