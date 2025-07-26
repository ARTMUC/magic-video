package customerauth

import (
	"context"

	"github.com/ARTMUC/magic-video/api/handler/customerauth/authdto"
	"github.com/ARTMUC/magic-video/internal/domain/customer"
	"github.com/ARTMUC/magic-video/internal/logger"
	"github.com/ARTMUC/magic-video/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type CustomerAuthController struct {
	customerService customer.CustomerService
	sessionService  service.SessionService
}

func NewCustomerAuthController(
	customerService customer.CustomerService,
	sessionService service.SessionService,
) *CustomerAuthController {
	return &CustomerAuthController{
		customerService: customerService,
		sessionService:  sessionService,
	}
}

func (c *CustomerAuthController) CreateAccess(
	ctx context.Context,
	input *authdto.CreateAccessRequestInput,
) (*struct{}, error) {
	// @TODO need to limit sending emails -> save email in db as one email not as an array of emails
	// and before sending email check if it was not send already this day (check email + template)
	_, err := c.customerService.CreateAccessThruEmail(input.Body.Email)
	if err != nil {
		logger.Log.Error("Error sending access email", zap.Error(err))
		return nil, huma.Error500InternalServerError("Internal Server Error")
	}

	return &struct{}{}, nil
}

func (c *CustomerAuthController) Signin(
	ctx context.Context,
	input *authdto.SigninRequestInput,
) (*authdto.CustomerAuthSigninOutput, error) {
	customer, err := c.customerService.GetCustomerFromToken(input.Body.CustomerUUID, input.Body.Token)
	if err != nil {
		logger.Log.Error("Error getting customer from access token", zap.Error(err))
		return nil, huma.Error401Unauthorized("Not authorized to sign in")
	}
	session, err := c.sessionService.CreateCustomerSession(customer)
	if err != nil {
		logger.Log.Error("Error creating customer session", zap.Error(err))
		return nil, huma.Error401Unauthorized("Not authorized to sign in")
	}

	return &authdto.CustomerAuthSigninOutput{
		Body: (&authdto.DtoConverterImpl{}).SessionToCustomerAuthSigninOutputBody(session),
	}, nil
}

func (c *CustomerAuthController) GetCustomer(
	ctx context.Context,
	input *struct{},
) (*authdto.GetCustomerOutput, error) {
	session, ok := c.sessionService.CustomerClaimsFromContext(ctx)
	if !ok {
		logger.Log.Error("Session not found in context")
		return nil, huma.Error400BadRequest("Invalid customer session")
	}

	return &authdto.GetCustomerOutput{
		Body: (&authdto.DtoConverterImpl{}).CustomerToGetCustomerOutputBody(session.Entity),
	}, nil
}
