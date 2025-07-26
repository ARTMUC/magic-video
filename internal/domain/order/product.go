package order

import (
	"fmt"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	base.BaseModel

	ProductTypeID uuid.UUID
	ProductType   *ProductType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Name      string          `gorm:"type:varchar(100)"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);type:TEXT"`
	TaxRate   decimal.Decimal `gorm:"type:decimal(5,2);type:TEXT"`
}

func (m *Product) BeforeUpdate(tx *gorm.DB) error {
	return fmt.Errorf("updates are not allowed; this model is immutable")
}

func (m *Product) BeforeSave(tx *gorm.DB) error {
	if tx.Statement.Changed() {
		return fmt.Errorf("updates are not allowed; this model is immutable")
	}
	return nil
}
