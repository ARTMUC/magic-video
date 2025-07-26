package customer

import (
	"time"

	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type CustomerAccess struct {
	base.BaseModel

	CustomerID      uuid.UUID
	Customer        *Customer
	AccessToken     string
	TokenExpireDate time.Time
}
