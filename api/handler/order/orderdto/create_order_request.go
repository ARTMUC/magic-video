package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	Body CreateOrderRequestBody
}

type CreateOrderRequestBody struct {
	Product    uuid.UUID   `json:"product" required:"true" format:"uuid" doc:"Product ID" example:"123e4567-e89b-12d3-a456-426614174000"`
	TemplateID uuid.UUID   `json:"template" required:"true" format:"uuid" doc:"Template ID selected by the customer" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
	PhotoIDs   []uuid.UUID `json:"photos" required:"true" minItems:"1" format:"uuid" doc:"List of photo IDs selected by the customer" example:"[550e8400-e29b-41d4-a716-446655440000, 6ba7b810-9dad-11d1-80b4-00c04fd430c8]"`
}
