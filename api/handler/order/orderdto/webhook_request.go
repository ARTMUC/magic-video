package orderdto

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type WebhookRequest struct {
	Request *http.Request
}

func (m *WebhookRequest) Resolve(ctx huma.Context) []error {
	m.Request = ctx.Context().Value("http_request").(*http.Request)

	return nil
}
