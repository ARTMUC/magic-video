package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BaseRepository[domain.Order]
}

type orderRepository struct {
	*BaseRepo[domain.Order]

	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepository{NewBaseRepository[domain.Order](db), db}
}

type OrderScopes struct {
	BaseScopes
}
