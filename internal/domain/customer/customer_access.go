package domain

import "time"

type CustomerAccess struct {
	BaseModel

	CustomerID      uint
	Customer        *Customer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT""`
	AccessToken     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	TokenExpireDate time.Time `gorm:"type:datetime;not null"`
}
