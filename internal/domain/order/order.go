package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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
	BaseModel

	CustomerID     uint
	Customer       *Customer       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	GrossAmount    decimal.Decimal `gorm:"type:decimal(10,2)"`
	NetAmount      decimal.Decimal `gorm:"type:decimal(10,2)"`
	TaxAmount      decimal.Decimal `gorm:"type:decimal(10,2)"`
	TaxBreakdown   TaxBreakdown    `gorm:"type:TEXT"`
	Status         string          `gorm:"type:varchar(255);not null"`
	PaymentStatus  string          `gorm:"type:varchar(255);not null"`
	IdempotencyKey string          `gorm:"type:varchar(255);not null;uniqueIndex"`

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
