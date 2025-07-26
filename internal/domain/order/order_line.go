package domain

import "github.com/shopspring/decimal"

type OrderLine struct {
	BaseModel

	OrderID            uint
	Order              *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VideoCompositionID uint
	VideoComposition   *VideoComposition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProductID          uint
	Product            *Product        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity           int             `gorm:"type:int;not null"`
	Amount             decimal.Decimal `gorm:"type:decimal(20,8);not null"`
}
