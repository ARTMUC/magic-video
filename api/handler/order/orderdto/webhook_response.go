package orderdto

type WebhookResponse struct {
	Body WebhookResponseBody
}

type WebhookResponseBody struct {
	Status string `json:"status" example:"OK"`
}
