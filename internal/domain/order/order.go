package order

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusCancelled = "cancelled"
	OrderStatusCompleted = "completed"
	OrderStatusFailed    = "failed"
	OrderStatusRefunded  = "refunded"
)

type Order struct {
	base.BaseModel

	CustomerID     uuid.UUID
	GrossAmount    decimal.Decimal
	NetAmount      decimal.Decimal
	TaxAmount      decimal.Decimal
	TaxBreakdown   TaxBreakdown
	Status         string
	PaymentStatus  string
	IdempotencyKey string

	OrderLines        []OrderLine
	OrderTransactions []OrderTransaction
	OrderPayments     []OrderPayment

	// @TODO
	//OrderRefunds []OrderRefund
}

type TaxBreakdown map[string]decimal.Decimal

func (t *TaxBreakdown) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into TaxBreakdown", value)
	}

	return json.Unmarshal(bytes, t)
}

func (t *TaxBreakdown) Value() (driver.Value, error) {
	return json.Marshal(t)
}
