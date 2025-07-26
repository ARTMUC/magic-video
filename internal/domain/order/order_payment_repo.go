package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type OrderPaymentRepository interface {
	BaseRepository[domain.OrderPayment]
}

type orderPaymentRepository struct {
	*BaseRepo[domain.OrderPayment]

	db *gorm.DB
}

func NewOrderPaymentRepo(db *gorm.DB) OrderPaymentRepository {
	return &orderPaymentRepository{NewBaseRepository[domain.OrderPayment](db), db}
}

type OrderPaymentScopes struct {
	BaseScopes
}

func (s OrderPaymentScopes) WithSessionID(sessionID string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("session_id = ?", sessionID)
	}
}
