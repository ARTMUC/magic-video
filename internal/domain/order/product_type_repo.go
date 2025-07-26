package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type ProductTypeRepository interface {
	base.BaseRepository[ProductType]
}

type productTypeRepository struct {
	*base.BaseRepo[ProductType]

	db *gorm.DB
}

func NewProductTypeRepo(db *gorm.DB) ProductTypeRepository {
	return &productTypeRepository{base.NewBaseRepository[ProductType](db), db}
}

type ProductTypeScopes struct {
	base.BaseScopes
}

const (
	ProductTypePreloadProduct = "Product"
)
