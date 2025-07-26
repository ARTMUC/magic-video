package customer

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	base.BaseRepository[Customer]
}

type customerRepo struct {
	*base.BaseRepo[Customer]

	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) CustomerRepository {
	return &customerRepo{base.NewBaseRepository[Customer](db), db}
}

type CustomerScopes struct {
	base.BaseScopes
}

func (u CustomerScopes) WithEmail(email string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customers.email = ?", email)
	}
}
