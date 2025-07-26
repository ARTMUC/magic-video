package job

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type VideoCompositionJobRepository interface {
	base.BaseRepository[VideoCompositionJob]
}

type videoCompositionJobRepository struct {
	*base.BaseRepo[VideoCompositionJob]

	db *gorm.DB
}

func NewVideoCompositionJobRepository(db *gorm.DB) VideoCompositionJobRepository {
	return &videoCompositionJobRepository{db: db}
}

type VideoCompositionJobScopes struct {
	base.BaseScopes
}

func (s VideoCompositionJobScopes) WithStatus(status string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}
