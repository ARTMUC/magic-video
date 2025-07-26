package crud

import (
	"github.com/ARTMUC/magic-video/internal/domain/customer"
)

type CustomerCrud interface {
	BaseCrud[customer.Customer]
}

type customerCrud struct {
	BaseCrud[customer.Customer]
	repository customer.CustomerRepository
}

func NewCustomerCrud(repository customer.CustomerRepository) CustomerCrud {
	return &customerCrud{BaseCrud: newBaseCrud(repository), repository: repository}
}
