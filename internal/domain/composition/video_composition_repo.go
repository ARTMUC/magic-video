package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type VideoCompositionRepository interface {
	BaseRepository[domain.VideoComposition]
}

type videoCompositionRepository struct {
	*BaseRepo[domain.VideoComposition]

	db *gorm.DB
}

func NewVideoCompositionRepository(db *gorm.DB) VideoCompositionRepository {
	return &videoCompositionRepository{db: db}
}

type VideoCompositionScopes struct {
	BaseScopes
}
