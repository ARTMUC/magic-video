package composition

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/order"
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

	VideoCompositionID uint
	VideoComposition   *VideoComposition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderLineID        uint
	OrderLine          *order.OrderLine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Status             string           `gorm:"type:varchar(255)"`
}
