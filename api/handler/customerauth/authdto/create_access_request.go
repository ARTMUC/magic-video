package authdto

type CreateAccessRequestInput struct {
	Body CreateAccessRequestInputBody
}

type CreateAccessRequestInputBody struct {
	Email string `json:"email" required:"true" format:"email"`
}
