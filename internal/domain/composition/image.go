package composition

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type Image struct {
	base.BaseModel

	Name               string
	PresetImageType    string
	VideoCompositionID uuid.UUID
}
