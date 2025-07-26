package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderLine struct {
	base.BaseModel

	OrderID            uuid.UUID
	VideoCompositionID uuid.UUID
	ProductID          uuid.UUID
	Quantity           int
	Amount             decimal.Decimal
}
