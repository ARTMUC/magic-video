package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type ActiveProductRepository interface {
	base.BaseRepository[ActiveProduct]
}

type activeProductRepository struct {
	*base.BaseRepo[ActiveProduct]

	db *gorm.DB
}

func NewProductTypeRepo(db *gorm.DB) ActiveProductRepository {
	return &activeProductRepository{base.NewBaseRepository[ActiveProduct](db), db}
}

type ActiveProductScopes struct {
	base.BaseScopes
}

const (
	ProductTypePreloadProduct = "Product"
)
