package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type OrderTransactionRepository interface {
	base.BaseRepository[OrderTransaction]
}

type orderTransactionRepository struct {
	*base.BaseRepo[OrderTransaction]

	db *gorm.DB
}

func NewOrderTransactionRepo(db *gorm.DB) OrderTransactionRepository {
	return &orderTransactionRepository{base.NewBaseRepository[OrderTransaction](db), db}
}

type OrderTransactionScopes struct {
	base.BaseScopes
}

func (s OrderTransactionScopes) WithSessionID(sessionID string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("session_id = ?", sessionID)
	}
}
