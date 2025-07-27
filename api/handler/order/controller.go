package order

import (
	"context"
	"errors"

	"github.com/ARTMUC/magic-video/api/handler/order/orderdto"
	"github.com/ARTMUC/magic-video/internal/domain/customer"
	"github.com/ARTMUC/magic-video/internal/domain/order"
	job "github.com/ARTMUC/magic-video/internal/domain/videojob"
	"github.com/ARTMUC/magic-video/internal/logger"
	"github.com/ARTMUC/magic-video/internal/service"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type OrderController struct {
	customerService customer.CustomerService
	sessionService  service.SessionService
	orderService    order.OrderService
	paymentService  order.PaymentService
	videoJobService job.VideoJobService
}

func NewOrderController(
	customerService customer.CustomerService,
	sessionService service.SessionService,
	orderService order.OrderService,
	paymentService order.PaymentService,
	videoJobService job.VideoJobService,
) *OrderController {
	return &OrderController{
		customerService: customerService,
		sessionService:  sessionService,
		orderService:    orderService,
		paymentService:  paymentService,
		videoJobService: videoJobService,
	}
}

func (c *OrderController) ProcessWebhook(
	ctx context.Context,
	input *orderdto.WebhookRequest,
) (*orderdto.WebhookResponse, error) {
	orderResult, err := c.paymentService.ProcessWebhook(input.Request)
	if err != nil {
		logger.Log.Error("Error processing webhook", zap.Error(err))
		if errors.Is(err, order.ErrPaymentServiceOrderNotFound) || errors.Is(err, order.ErrPaymentServiceOrderTransactionNotFound) {
			return nil, huma.Error400BadRequest("Invalid sessionId")
		}
	}

	err = c.videoJobService.Enqueue(orderResult.ID)
	if err != nil {
		logger.Log.Error("Error enqueuing video composition job", zap.Error(err))
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &orderdto.WebhookResponse{
		Body: orderdto.WebhookResponseBody{
			Status: "OK",
		},
	}, nil
}

func (c *OrderController) CreateOrder(
	ctx context.Context,
	input *orderdto.CreateOrderRequest,
) (*orderdto.CreateOrderResponse, error) {
	session, ok := c.sessionService.CustomerClaimsFromContext(ctx)
	if !ok {
		logger.Log.Error("Session not found in context")
		return nil, huma.Error403Forbidden("Invalid customer session")
	}

	orderResult, err := c.orderService.ProcessCart(session.Entity, input.Body.Transform(&order.CreateOrderInput{}))
	if err != nil {
		logger.Log.Error("Error processing cart", zap.Error(err))
		switch {
		case errors.Is(err, order.ErrOrderServiceProductNotFound),
			errors.Is(err, order.ErrOrderServiceVideoCompositionNotFound):
			return nil, huma.Error400BadRequest("Invalid request body")
		}
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	transaction, err := c.paymentService.CreateTransaction(orderResult, session.Entity)
	if err != nil {
		logger.Log.Error("Error creating payment", zap.Error(err))
		return nil, huma.Error500InternalServerError("Internal server error")
	}

	return &orderdto.CreateOrderResponse{
		Body: &orderdto.CreateOrderResponseBody{
			SessionID:  transaction.SessionIden,
			PaymentURL: transaction.PaymentUrl,
			Token:      transaction.Token,
		},
	}, nil
}
