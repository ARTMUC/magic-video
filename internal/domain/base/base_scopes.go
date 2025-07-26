package base

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseScopes struct {
}

func (u BaseScopes) WithID(id uuid.UUID) Scope {
	return WithID(id)
}

func (u BaseScopes) OrderBy(clause string) Scope {
	return OrderBy(clause)
}

func WithID(id uuid.UUID) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func OrderBy(clause string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause)
	}
}
