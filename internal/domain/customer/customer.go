package domain

import "database/sql"

type Customer struct {
	BaseModel

	Email              string           `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name               sql.Null[string] `gorm:"type:varchar(255)"`
	AddressHomeNumber  sql.Null[string] `gorm:"type:varchar(255)"`
	AddressStreetName  sql.Null[string] `gorm:"type:varchar(255)"`
	AddressCity        sql.Null[string] `gorm:"type:varchar(255)"`
	AddressZipCode     sql.Null[string] `gorm:"type:varchar(255)"`
	AddressState       sql.Null[string] `gorm:"type:varchar(255)"`
	AddressCountryCode sql.Null[string] `gorm:"type:varchar(255)"`
}
