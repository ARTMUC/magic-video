package contracts

import (
	"database/sql"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func BaseToBase(source base.BaseModel) BaseModel {
	var contractsBaseModel BaseModel
	contractsBaseModel.ID = UUIDToUUID(source.ID)
	return contractsBaseModel
}

func UUIDToUUID(source uuid.UUID) uuid.UUID {
	return source
}

func DecimalToDecimal(source decimal.Decimal) decimal.Decimal {
	return source
}

func NullStringToStringPtr(source sql.Null[string]) *string {
	if !source.Valid {
		return nil
	}
	return &source.V
}
