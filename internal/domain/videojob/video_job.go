package job

import (
	"database/sql"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

const (
	VideoCompositionJobStatusCompleted  = "completed"
	VideoCompositionJobStatusInProgress = "in_progress"
	VideoCompositionJobStatusFailed     = "failed"
	VideoCompositionJobStatusCancelled  = "cancelled"
	VideoCompositionJobStatusCreated    = "created"
)

type VideoCompositionJob struct {
	base.BaseModel

	OrderID uuid.UUID
	Status  string
	Error   sql.Null[string]
}
