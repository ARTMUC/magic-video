package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:varchar(255);not null;uniqueIndex"`
}

func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	return
}
