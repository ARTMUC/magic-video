package orderdto

import (
	"github.com/ARTMUC/magic-video/internal/domain/order"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Body CreateOrderRequestBody
}

type CreateOrderRequestBody struct {
	Cart           []CartItem `json:"cart" doc:"List of items in the cart" minItems:"1" maxItems:"100"`
	IdempotencyKey uuid.UUID  `json:"idempotency_key" doc:"Unique key to prevent duplicate orders" format:"uuid"`
}

func (in CreateOrderRequestBody) Transform(out *order.CreateOrderInput) *order.CreateOrderInput {
	out.IdempotencyKey = in.IdempotencyKey
	out.Cart = make([]*order.CreateOrderCartItemInput, len(in.Cart))
	for i, item := range in.Cart {
		out.Cart[i] = item.Transform(&order.CreateOrderCartItemInput{})
	}

	return out
}

type CartItem struct {
	ProductTypeUUID      uuid.UUID `json:"product_uuid" doc:"Product's type unique identifier" format:"uuid"`
	VideoCompositionUUID uuid.UUID `json:"video_composition_uuid" doc:"Video composition's unique identifier" format:"uuid"`
	Quantity             int       `json:"quantity" doc:"Number of items" minimum:"1" maximum:"1000"`
}

func (in CartItem) Transform(out *order.CreateOrderCartItemInput) *order.CreateOrderCartItemInput {
	out.ProductTypeUUID = in.ProductTypeUUID
	out.VideoCompositionUUID = in.VideoCompositionUUID
	out.Quantity = in.Quantity

	return out
}
