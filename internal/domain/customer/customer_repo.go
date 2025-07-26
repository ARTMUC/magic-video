package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	BaseRepository[domain.Customer]
}

type customerRepo struct {
	*BaseRepo[domain.Customer]

	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) CustomerRepository {
	return &customerRepo{NewBaseRepository[domain.Customer](db), db}
}

type CustomerScopes struct {
	BaseScopes
}

func (u CustomerScopes) WithEmail(email string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customers.email = ?", email)
	}
}
