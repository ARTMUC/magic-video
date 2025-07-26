package composition

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type ImageRepository interface {
	base.BaseRepository[Image]
}

type imageRepository struct {
	*base.BaseRepo[Image]

	db *gorm.DB
}

func NewImageRepo(db *gorm.DB) ImageRepository {
	return &imageRepository{base.NewBaseRepository[Image](db), db}
}

type ImageScopes struct {
	base.BaseScopes
}
