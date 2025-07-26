package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type OrderTransactionRepository interface {
	BaseRepository[domain.OrderTransaction]
}

type orderTransactionRepository struct {
	*BaseRepo[domain.OrderTransaction]

	db *gorm.DB
}

func NewOrderTransactionRepo(db *gorm.DB) OrderTransactionRepository {
	return &orderTransactionRepository{NewBaseRepository[domain.OrderTransaction](db), db}
}

type OrderTransactionScopes struct {
	BaseScopes
}

func (s OrderTransactionScopes) WithSessionID(sessionID string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("session_id = ?", sessionID)
	}
}
