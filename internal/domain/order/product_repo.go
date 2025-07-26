package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[domain.Product]
}

type productRepository struct {
	*BaseRepo[domain.Product]

	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &productRepository{NewBaseRepository[domain.Product](db), db}
}

type ProductScopes struct {
	BaseScopes
}
