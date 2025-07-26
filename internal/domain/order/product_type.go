package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type ProductType struct {
	base.BaseModel

	ProductID uuid.UUID // current product version
	Product   *Product  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
