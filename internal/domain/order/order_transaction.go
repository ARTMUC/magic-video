package domain

type OrderTransaction struct {
	BaseModel

	OrderID uint
	Order   *Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	
	Amount          int    `gorm:"type:int"`
	Method          string `gorm:"type:varchar(100)"`
	Token           string `gorm:"type:varchar(255)"`
	SessionIden     string `gorm:"unique_index"`
	TransactionIden string `gorm:"unique_index"`
	PaymentUrl      string `gorm:"type:text"`
}
