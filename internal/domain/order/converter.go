//go:generate goverter gen .

package order

import (
	"github.com/ARTMUC/magic-video/internal/contracts"
)

var BaseToBase = contracts.BaseToBase
var UUIDToUUID = contracts.UUIDToUUID
var DecimalToDecimal = contracts.DecimalToDecimal

// goverter:converter
// goverter:extend BaseToBase
// goverter:extend UUIDToUUID
// goverter:extend DecimalToDecimal
// goverter:output:file ./converter_generated.go
type VideoCompositionConverter interface {
	VideoCompositionDomainToContract(source Order) contracts.Order
}
