package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type OrderTransaction struct {
	base.BaseModel

	OrderID uuid.UUID
	Order   *Order

	Amount          int
	Method          string
	Token           string
	SessionIden     string
	TransactionIden string
	PaymentUrl      string
}
