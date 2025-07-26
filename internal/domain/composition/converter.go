//go:generate goverter gen .

package composition

import (
	"github.com/ARTMUC/magic-video/internal/contracts"
)

var BaseToBase = contracts.BaseToBase
var UUIDToUUID = contracts.UUIDToUUID

// goverter:converter
// goverter:extend BaseToBase
// goverter:extend UUIDToUUID
// goverter:output:file ./converter_generated.go
type VideoCompositionConverter interface {
	VideoCompositionDomainToContract(source VideoComposition) contracts.VideoComposition
	VideoCompositionDomainToContractArray(source []VideoComposition) []contracts.VideoComposition
	VideoCompositionImageToContract(source Image) contracts.Image
}
