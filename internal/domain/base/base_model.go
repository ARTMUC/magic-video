package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
