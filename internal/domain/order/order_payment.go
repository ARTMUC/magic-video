package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

const (
	OrderPaymentStatusPending   = "pending"
	OrderPaymentStatusCancelled = "cancelled"
	OrderPaymentStatusCompleted = "completed"
	OrderPaymentStatusFailed    = "failed"
	OrderPaymentStatusRefunded  = "refunded"
)

type OrderPayment struct {
	base.BaseModel

	OrderTransactionID uuid.UUID
	OrderTransaction   *OrderTransaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderID            uuid.UUID
	Order              *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SessionID          string `gorm:"unique_index"`
}
