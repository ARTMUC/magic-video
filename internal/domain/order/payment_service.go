package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain"
	"github.com/ARTMUC/magic-video/internal/pkg/p24"
	"github.com/ARTMUC/magic-video/internal/repository"
)

var (
	ErrPaymentServiceOrderTransactionNotFound = errors.New("payment service order transaction not found")
	ErrPaymentServiceOrderNotFound            = errors.New("payment service order not found")
)

type PaymentService interface {
	CreateTransaction(order *domain.Order, customer *domain.Customer) (*domain.OrderTransaction, error)
	ProcessWebhook(r *http.Request) (*domain.Order, error)
}

type paymentService struct {
	paymentConfig              config.Przelewy24ClientConfig
	transactionProvider        repository.TransactionProvider
	orderTransactionRepository repository.OrderTransactionRepository
	orderPaymentRepository     repository.OrderPaymentRepository
	orderRepository            repository.OrderRepository
}

func NewPaymentService(
	paymentConfig config.Przelewy24ClientConfig,
	transactionProvider repository.TransactionProvider,
	orderTransactionRepository repository.OrderTransactionRepository,
	orderPaymentRepository repository.OrderPaymentRepository,
	orderRepository repository.OrderRepository,
) PaymentService {
	return &paymentService{
		paymentConfig:              paymentConfig,
		transactionProvider:        transactionProvider,
		orderTransactionRepository: orderTransactionRepository,
		orderPaymentRepository:     orderPaymentRepository,
		orderRepository:            orderRepository,
	}
}

func (s *paymentService) getClient() *p24.Client {
	return p24.NewClient(
		s.paymentConfig.MerchatID(),
		s.paymentConfig.PosID(),
		s.paymentConfig.Salt(),
		s.paymentConfig.ApiKey(),
		s.paymentConfig.Environment(),
	)
}

func (s *paymentService) CreateTransaction(
	order *domain.Order,
	customer *domain.Customer,
) (*domain.OrderTransaction, error) {
	client := s.getClient()
	sessionID := fmt.Sprintf("order_%d_%d", order.UUID, time.Now().Nanosecond())
	transaction := p24.NewTransactionBuilder().
		SetSessionID(sessionID).
		SetAmountDecimal(order.GrossAmount).
		SetEmail(customer.Email).
		SetReturnURL(s.paymentConfig.ReturnUrl()).
		SetStatusURL(s.paymentConfig.WebhookUrl()).
		Build()

	response, err := client.RegisterTransaction(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to register transaction: %w", err)
	}

	var orderTransaction *domain.OrderTransaction
	err = s.transactionProvider.Transaction(func(tx *repository.Tx) error {
		orderTransaction = &domain.OrderTransaction{
			Amount:      transaction.Amount,
			Method:      "online",
			SessionIden: sessionID,
			Token:       response.Data.Token,
			PaymentUrl:  client.GetPaymentURL(response.Data.Token),
			OrderID:     order.ID,
		}

		err = s.orderTransactionRepository.Create(repository.WriteOptions{Tx: tx}, orderTransaction)
		if err != nil {
			return fmt.Errorf("failed to create order transaction in db: %w", err)
		}

		order.PaymentStatus = domain.OrderPaymentStatusPending
		err = s.orderRepository.Update(repository.WriteOptions{Tx: tx}, order)
		if err != nil {
			return fmt.Errorf("failed to update order in db: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return orderTransaction, nil
}

func (s *paymentService) ProcessWebhook(r *http.Request) (*domain.Order, error) {
	client := s.getClient()
	data, err := p24.ParseWebhookData(r)
	if err != nil {
		return nil, fmt.Errorf("can't parse webhook data: %w", err)
	}

	ok := client.VerifyWebhookSignature(data)
	if !ok {
		return nil, fmt.Errorf("invalid signature")
	}

	orderTransaction, err := s.orderTransactionRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{repository.OrderTransactionScopes{}.WithSessionID(data.SessionId)},
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrPaymentServiceOrderTransactionNotFound
		}
		return nil, fmt.Errorf("can't find order transaction in db: %w", err)
	}
	order, err := s.orderRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{repository.OrderScopes{}.WithID(orderTransaction.OrderID)},
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrPaymentServiceOrderNotFound
		}
		return nil, fmt.Errorf("can't find order in db: %w", err)
	}

	var orderPayment *domain.OrderPayment
	err = s.transactionProvider.Transaction(func(tx *repository.Tx) error {
		orderPayment = &domain.OrderPayment{
			BaseModel:          domain.BaseModel{},
			OrderTransactionID: orderTransaction.ID,
			OrderID:            order.ID,
			SessionID:          data.SessionId,
		}

		err = s.orderPaymentRepository.Create(repository.WriteOptions{Tx: tx}, orderPayment)
		if err != nil {
			return fmt.Errorf("failed to create order payment in db: %w", err)
		}

		order.PaymentStatus = domain.OrderPaymentStatusCompleted
		err = s.orderRepository.Update(repository.WriteOptions{Tx: tx}, order)
		if err != nil {
			return fmt.Errorf("failed to update order in db: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}
