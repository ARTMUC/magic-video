package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLineRepository interface {
	base.BaseRepository[OrderLine]
}

type orderLineRepository struct {
	*base.BaseRepo[OrderLine]

	db *gorm.DB
}

func NewOrderLineRepo(db *gorm.DB) OrderLineRepository {
	return &orderLineRepository{base.NewBaseRepository[OrderLine](db), db}
}

type OrderLineScopes struct {
	base.BaseScopes
}

const (
	OrderLinePreloadVideoComposition = "VideoComposition"
)

func (o OrderLineScopes) WithOrderID(orderID uuid.UUID) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("order_id = ?", orderID)
	}
}
