package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type OrderRepository interface {
	base.BaseRepository[Order]
}

type orderRepository struct {
	*base.BaseRepo[Order]

	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepository{base.NewBaseRepository[Order](db), db}
}

type OrderScopes struct {
	base.BaseScopes
}
