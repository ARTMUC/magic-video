package customerauth

import (
	"net/http"

	"github.com/ARTMUC/magic-video/api/middleware"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, c *CustomerAuthController) {
	huma.Register(api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/customers/access",
			DefaultStatus: 201,
			Middlewares:   huma.Middlewares{},

			Tags:        []string{"CustomersAuth"},
			Summary:     "Create an access for customer",
			Description: "Create an access for customer and sends an access email address to the customer.",
			OperationID: "createCustomerAccess",
		},
		c.CreateAccess,
	)
	huma.Register(api,
		huma.Operation{
			Method:        http.MethodPost,
			Path:          "/customers/session",
			DefaultStatus: 201,
			Middlewares:   huma.Middlewares{},

			Tags:        []string{"CustomersAuth"},
			Summary:     "Creates session for customer",
			Description: "Create session for customer.",
			OperationID: "createCustomerSession",
		},
		c.Signin,
	)
	huma.Register(api,
		huma.Operation{
			Method:        http.MethodGet,
			Path:          "/customers/session",
			DefaultStatus: 200,
			Middlewares: huma.Middlewares{
				middleware.Auth(api, c.sessionService.ParseCustomerToken),
			},
			Tags:        []string{"CustomersAuth"},
			Summary:     "Returns customer",
			Description: "Returns customer.",
			OperationID: "getCustomerSession",
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
		},
		c.GetCustomer,
	)
}
