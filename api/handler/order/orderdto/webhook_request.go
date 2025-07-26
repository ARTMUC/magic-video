package dto

type WebhookRequest struct {
	Body map[string]string `json:"-"` // We'll parse form data manually
}
