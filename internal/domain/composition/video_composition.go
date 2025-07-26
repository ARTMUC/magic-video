package composition

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

const (
	VideoCompositionStatusCompleted  = "completed"
	VideoCompositionStatusInProgress = "in_progress"
	VideoCompositionStatusFailed     = "failed"
	VideoCompositionStatusCancelled  = "cancelled"
	VideoCompositionStatusCreated    = "created"
)

type VideoComposition struct {
	base.BaseModel

	CustomerID    uuid.UUID
	VideoTemplate string
	Status        string
	Images        []Image
}
