package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type OrderLineRepository interface {
	BaseRepository[domain.OrderLine]
}

type orderLineRepository struct {
	*BaseRepo[domain.OrderLine]

	db *gorm.DB
}

func NewOrderLineRepo(db *gorm.DB) OrderLineRepository {
	return &orderLineRepository{NewBaseRepository[domain.OrderLine](db), db}
}

type OrderLineScopes struct {
	BaseScopes
}

const (
	OrderLinePreloadVideoComposition = "VideoComposition"
)

func (o OrderLineScopes) WithOrderID(orderID uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}
