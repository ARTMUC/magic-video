package domain

type ProductType struct {
	BaseModel

	ProductID uint     // current product version
	Product   *Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
