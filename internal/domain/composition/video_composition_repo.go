package composition

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type VideoCompositionRepository interface {
	base.BaseRepository[VideoComposition]
}

type videoCompositionRepository struct {
	*base.BaseRepo[VideoComposition]

	db *gorm.DB
}

func NewVideoCompositionRepository(db *gorm.DB) VideoCompositionRepository {
	return &videoCompositionRepository{db: db}
}

type VideoCompositionScopes struct {
	base.BaseScopes
}
