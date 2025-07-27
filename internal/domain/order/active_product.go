package order

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/google/uuid"
)

type ActiveProduct struct {
	base.BaseModel

	ProductID uuid.UUID // current product version
	Product   *Product
}
