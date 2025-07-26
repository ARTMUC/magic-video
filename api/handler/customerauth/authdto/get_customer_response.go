package authdto

import (
	"github.com/ARTMUC/magic-video/internal/contracts"
)

type GetCustomerOutput struct {
	Body *GetCustomerOutputBody
}

// TGen customer.Customer [reverse]
type GetCustomerOutputBody struct {
	contracts.BaseModel

	Email              string  `json:"email,omitempty"`
	Name               *string `json:"name,omitempty"`
	AddressHomeNumber  *string `json:"addressHomeNumber,omitempty"`
	AddressStreetName  *string `json:"addressStreetName,omitempty"`
	AddressCity        *string `json:"addressCity,omitempty"`
	AddressZipCode     *string `json:"addressZipCode,omitempty"`
	AddressState       *string `json:"addressState,omitempty"`
	AddressCountryCode *string `json:"addressCountryCode,omitempty"`
}
