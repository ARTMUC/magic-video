package repository

import (
	"time"

	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type CustomerAccessRepository interface {
	BaseRepository[domain.CustomerAccess]
}

type customerAccessRepo struct {
	*BaseRepo[domain.CustomerAccess]

	db *gorm.DB
}

func NewCustomerAccessRepo(db *gorm.DB) CustomerAccessRepository {
	return &customerAccessRepo{NewBaseRepository[domain.CustomerAccess](db), db}
}

type CustomerAccessScopes struct {
	BaseScopes
}

func (u CustomerAccessScopes) WithCustomer(customer *domain.Customer) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_accesses.customer_id = ?", customer.ID)
	}
}

func (u CustomerAccessScopes) WithNotExpired() Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_accesses.token_expire_date > ?", time.Now())
	}
}
