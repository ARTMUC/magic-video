package customer

import (
	"time"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type CustomerAccessRepository interface {
	base.BaseRepository[CustomerAccess]
}

type customerAccessRepo struct {
	*base.BaseRepo[CustomerAccess]

	db *gorm.DB
}

func NewCustomerAccessRepo(db *gorm.DB) CustomerAccessRepository {
	return &customerAccessRepo{base.NewBaseRepository[CustomerAccess](db), db}
}

type CustomerAccessScopes struct {
	base.BaseScopes
}

func (u CustomerAccessScopes) WithCustomer(customer *Customer) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_accesses.customer_id = ?", customer.ID)
	}
}

func (u CustomerAccessScopes) WithNotExpired() base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("customer_accesses.token_expire_date > ?", time.Now())
	}
}
