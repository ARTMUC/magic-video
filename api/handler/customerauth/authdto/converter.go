//go:generate goverter gen .

package authdto

import (
	"github.com/ARTMUC/magic-video/internal/contracts"
	"github.com/ARTMUC/magic-video/internal/domain/customer"
	"github.com/ARTMUC/magic-video/internal/service"
)

var BaseToBase = contracts.BaseToBase
var UUIDToUUID = contracts.UUIDToUUID
var NullStringToStringPtr = contracts.NullStringToStringPtr

// goverter:converter
// goverter:extend BaseToBase
// goverter:extend UUIDToUUID
// goverter:extend NullStringToStringPtr
// goverter:output:file ./converter_generated.go
type DtoConverter interface {
	CustomerToGetCustomerOutputBody(source *customer.Customer) *GetCustomerOutputBody
	SessionToCustomerAuthSigninOutputBody(source *service.SessionOutput) *CustomerAuthSigninOutputBody
}
