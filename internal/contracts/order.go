package contracts

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	BaseModel

	CustomerID    uuid.UUID
	GrossAmount   decimal.Decimal
	NetAmount     decimal.Decimal
	TaxAmount     decimal.Decimal
	Status        string
	PaymentStatus string

	OrderLines []OrderLine
}

type OrderLine struct {
	BaseModel

	OrderID            uuid.UUID
	VideoCompositionID uuid.UUID
	ProductID          uuid.UUID
	Quantity           int
	Amount             decimal.Decimal
}
