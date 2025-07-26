package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type ProductRepository interface {
	base.BaseRepository[Product]
}

type productRepository struct {
	*base.BaseRepo[Product]

	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &productRepository{base.NewBaseRepository[Product](db), db}
}

type ProductScopes struct {
	base.BaseScopes
}
