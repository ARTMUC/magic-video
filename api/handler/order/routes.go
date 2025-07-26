package order

import (
	"net/http"

	"github.com/ARTMUC/magic-video/api/middleware"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, c *OrderController) {
	huma.Register(api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/payments/p24/notification",
			DefaultStatus: 201,
			Middlewares:   huma.Middlewares{},
			Tags:          []string{"Order"},
			Summary:       "Processes payment notification",
			Description:   "Processes payment notification",
			OperationID:   "processWebhook",
			Security:      []map[string][]string{},
		},
		c.ProcessWebhook,
	)
	huma.Register(api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/orders",
			DefaultStatus: 201,
			Middlewares: huma.Middlewares{
				middleware.Auth(api, c.sessionService.ParseCustomerToken),
			},
			Tags:        []string{"Order"},
			Summary:     "Creates order.",
			Description: "Creates order.",
			OperationID: "createOrder",
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		},
		c.CreateOrder,
	)
}
