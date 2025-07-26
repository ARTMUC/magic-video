package customer

import (
	"database/sql"

	"github.com/ARTMUC/magic-video/internal/domain/base"
)

type Customer struct {
	base.BaseModel

	Email              string
	Name               sql.Null[string]
	AddressHomeNumber  sql.Null[string]
	AddressStreetName  sql.Null[string]
	AddressCity        sql.Null[string]
	AddressZipCode     sql.Null[string]
	AddressState       sql.Null[string]
	AddressCountryCode sql.Null[string]
}
