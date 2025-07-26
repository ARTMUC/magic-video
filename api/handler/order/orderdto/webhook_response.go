package dto

type WebhookResponse struct {
	Body struct {
		Status  string `json:"status" example:"OK"`
		Message string `json:"message,omitempty" example:"Payment verified successfully"`
	}
}
