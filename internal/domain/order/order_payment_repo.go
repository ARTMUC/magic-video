package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type OrderPaymentRepository interface {
	base.BaseRepository[OrderPayment]
}

type orderPaymentRepository struct {
	*base.BaseRepo[OrderPayment]

	db *gorm.DB
}

func NewOrderPaymentRepo(db *gorm.DB) OrderPaymentRepository {
	return &orderPaymentRepository{base.NewBaseRepository[OrderPayment](db), db}
}

type OrderPaymentScopes struct {
	base.BaseScopes
}

func (s OrderPaymentScopes) WithSessionID(sessionID string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("session_id = ?", sessionID)
	}
}
