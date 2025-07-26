package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type ImageRepository interface {
	BaseRepository[domain.Image]
}

type imageRepository struct {
	*BaseRepo[domain.Image]

	db *gorm.DB
}

func NewImageRepo(db *gorm.DB) ImageRepository {
	return &imageRepository{NewBaseRepository[domain.Image](db), db}
}

type ImageScopes struct {
	BaseScopes
}
