package authdto

import "github.com/google/uuid"

type SigninRequestInput struct {
	Body SigninRequestInputBody
}

type SigninRequestInputBody struct {
	Token        string    `json:"token" required:"true"`
	CustomerUUID uuid.UUID `json:"customer" required:"true"`
}
