package domain

const (
	OrderPaymentStatusPending   = "pending"
	OrderPaymentStatusCancelled = "cancelled"
	OrderPaymentStatusCompleted = "completed"
	OrderPaymentStatusFailed    = "failed"
	OrderPaymentStatusRefunded  = "refunded"
)

type OrderPayment struct {
	BaseModel

	OrderTransactionID uint
	OrderTransaction   *OrderTransaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderID            uint
	Order              *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SessionID          string `gorm:"unique_index"`
}
