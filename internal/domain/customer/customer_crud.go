package customer

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
)

type CustomerCrud interface {
	base.BaseCrud[Customer]
}

type customerCrud struct {
	base.BaseCrud[Customer]
	repository CustomerRepository
}

func NewCustomerCrud(repository CustomerRepository) CustomerCrud {
	return &customerCrud{BaseCrud: base.NewBaseCrud(repository), repository: repository}
}
