package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseScopes struct {
}

func (u BaseScopes) WithID(id uint) Scope {
	return WithID(id)
}

func (u BaseScopes) WithUUID(uuid uuid.UUID) Scope {
	return WithUUID(uuid)
}

func (u BaseScopes) OrderBy(clause string) Scope {
	return OrderBy(clause)
}

func WithID(id uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func WithUUID(uuid uuid.UUID) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid = ?", uuid.String())
	}
}

func OrderBy(clause string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause)
	}
}
