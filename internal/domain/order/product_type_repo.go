package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type ProductTypeRepository interface {
	BaseRepository[domain.ProductType]
}

type productTypeRepository struct {
	*BaseRepo[domain.ProductType]

	db *gorm.DB
}

func NewProductTypeRepo(db *gorm.DB) ProductTypeRepository {
	return &productTypeRepository{NewBaseRepository[domain.ProductType](db), db}
}

type ProductTypeScopes struct {
	BaseScopes
}

const (
	ProductTypePreloadProduct = "Product"
)
