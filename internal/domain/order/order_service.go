package service

import (
	"errors"
	"fmt"

	"github.com/ARTMUC/magic-video/internal/domain"
	"github.com/ARTMUC/magic-video/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrOrderServiceVideoCompositionNotFound = errors.New("video composition not found")
	ErrOrderServiceProductNotFound          = errors.New("product not found")
)

type OrderService interface {
	ProcessCart(customer *domain.Customer, input *CreateOrderInput) (*domain.Order, error)
}

type orderService struct {
	transactionProvider        repository.TransactionProvider
	videoCompositionRepository repository.VideoCompositionRepository
	productTypeRepository      repository.ProductTypeRepository
	orderRepository            repository.OrderRepository
	orderLineRepository        repository.OrderLineRepository
}

func NewOrderService(
	transactionProvider repository.TransactionProvider,
	videoCompositionRepository repository.VideoCompositionRepository,
	productTypeRepository repository.ProductTypeRepository,
	orderRepository repository.OrderRepository,
	orderLineRepository repository.OrderLineRepository,
) OrderService {
	return &orderService{
		transactionProvider:        transactionProvider,
		videoCompositionRepository: videoCompositionRepository,
		productTypeRepository:      productTypeRepository,
		orderRepository:            orderRepository,
		orderLineRepository:        orderLineRepository,
	}
}

type CreateOrderInput struct {
	Cart           []*CreateOrderCartItemInput
	IdempotencyKey uuid.UUID
}

type CreateOrderCartItemInput struct {
	ProductTypeUUID      uuid.UUID
	VideoCompositionUUID uuid.UUID
	Quantity             int
}

type CartItemDetails struct {
	Product              *domain.Product
	VideoComposition     *domain.VideoComposition
	Quantity             int
	NetAmountUnrounded   decimal.Decimal
	VatAmountUnrounded   decimal.Decimal
	GrossAmountUnrounded decimal.Decimal
	NetAmountRounded     decimal.Decimal
	VATAmountRounded     decimal.Decimal
	GrossAmountRounded   decimal.Decimal
}

func (o *orderService) ProcessCart(customer *domain.Customer, input *CreateOrderInput) (*domain.Order, error) {
	productTypes, err := o.loadProducts(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	videoCompositions, err := o.loadVideoCompositions(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load video compositions: %w", err)
	}

	taxBreakdownMap := make(map[string]decimal.Decimal)
	cartItems := make([]CartItemDetails, len(input.Cart))
	totalNetUnrounded := decimal.Zero
	totalVATUnrounded := decimal.Zero
	totalGrossUnrounded := decimal.Zero

	for i, item := range input.Cart {
		productType := productTypes[item.ProductTypeUUID]
		videoComposition := videoCompositions[item.VideoCompositionUUID]

		if productType.Product == nil {
			return nil, fmt.Errorf("not loaded preload Product in ProductType")
		}

		product := productType.Product

		unitPrice := product.UnitPrice
		taxRate := product.TaxRate

		lineNetAmountUnrounded := unitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		lineVatAmountUnrounded := lineNetAmountUnrounded.Mul(taxRate.Div(decimal.NewFromInt(100)))
		lineGrossAmountUnrounded := lineNetAmountUnrounded.Add(lineVatAmountUnrounded)

		lineNetAmountRounded := lineNetAmountUnrounded.Round(2)
		lineVatAmountRounded := lineVatAmountUnrounded.Round(2)
		lineGrossAmountRounded := lineGrossAmountUnrounded.Round(2)

		cartItem := CartItemDetails{
			Product:              product,
			VideoComposition:     videoComposition,
			Quantity:             item.Quantity,
			NetAmountUnrounded:   lineNetAmountUnrounded,
			VatAmountUnrounded:   lineVatAmountUnrounded,
			GrossAmountUnrounded: lineGrossAmountUnrounded,
			NetAmountRounded:     lineNetAmountRounded,
			VATAmountRounded:     lineVatAmountRounded,
			GrossAmountRounded:   lineGrossAmountRounded,
		}

		cartItems[i] = cartItem

		totalNetUnrounded = totalNetUnrounded.Add(lineNetAmountUnrounded)
		totalVATUnrounded = totalVATUnrounded.Add(lineVatAmountUnrounded)
		totalGrossUnrounded = totalGrossUnrounded.Add(lineGrossAmountUnrounded)

		taxRateStr := product.TaxRate.String()
		taxBreakdownMap[taxRateStr] = taxBreakdownMap[taxRateStr].Add(lineVatAmountUnrounded)
	}

	taxBreakdownRounded := make(map[string]decimal.Decimal)
	for rate, amount := range taxBreakdownMap {
		taxBreakdownRounded[rate] = amount.Round(2)
	}

	var order *domain.Order
	err = o.transactionProvider.Transaction(func(tx *repository.Tx) error {
		order = &domain.Order{
			CustomerID:     customer.ID,
			Customer:       customer,
			GrossAmount:    totalGrossUnrounded.Round(2),
			NetAmount:      totalNetUnrounded.Round(2),
			TaxAmount:      totalVATUnrounded.Round(2),
			TaxBreakdown:   taxBreakdownRounded,
			Status:         domain.OrderStatusPending,
			PaymentStatus:  domain.OrderPaymentStatusPending,
			IdempotencyKey: fmt.Sprintf("%s:%s", customer.UUID, input.IdempotencyKey),
			OrderLines:     make([]domain.OrderLine, len(cartItems)),
		}

		err := o.orderRepository.Create(repository.WriteOptions{Tx: tx}, order)
		if err != nil {
			return fmt.Errorf("failed to create order in db: %w", err)
		}

		for i, cartItem := range cartItems {
			orderLine := &domain.OrderLine{
				Quantity:           cartItem.Quantity,
				OrderID:            order.ID,
				VideoCompositionID: cartItem.VideoComposition.ID,
				ProductID:          cartItem.Product.ID,
				Amount:             cartItem.GrossAmountRounded,
			}
			err = o.orderLineRepository.Create(repository.WriteOptions{Tx: tx}, orderLine)
			if err != nil {
				return fmt.Errorf("failed to create order line in db: %w", err)
			}
			order.OrderLines[i] = *orderLine
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (o *orderService) loadVideoCompositions(input *CreateOrderInput) (map[uuid.UUID]*domain.VideoComposition, error) {
	videoCompositions := make(map[uuid.UUID]*domain.VideoComposition)
	for _, item := range input.Cart {
		if videoCompositions[item.VideoCompositionUUID] != nil {
			continue
		}
		videoComposition, err := o.videoCompositionRepository.FindOne(repository.ReadOptions{
			Scopes: []repository.Scope{
				repository.VideoCompositionScopes{}.WithUUID(item.VideoCompositionUUID),
			},
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: video composition: %s not found: %w", ErrOrderServiceVideoCompositionNotFound, item.VideoCompositionUUID, err)
			}
			return nil, fmt.Errorf("failed to find video composition: %w", err)
		}
		videoCompositions[item.VideoCompositionUUID] = videoComposition
	}

	return videoCompositions, nil
}

// loadProducts: loads product types from the database with preloaded current product
func (o *orderService) loadProducts(input *CreateOrderInput) (map[uuid.UUID]*domain.ProductType, error) {
	productTypes := make(map[uuid.UUID]*domain.ProductType)
	for _, item := range input.Cart {
		if productTypes[item.ProductTypeUUID] != nil {
			continue
		}
		product, err := o.productTypeRepository.FindOne(repository.ReadOptions{
			Scopes: []repository.Scope{
				repository.ProductTypeScopes{}.WithUUID(item.ProductTypeUUID),
			},
			Preload: []string{repository.ProductTypePreloadProduct},
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: product type: %s not found: %w", ErrOrderServiceProductNotFound, item.ProductTypeUUID, err)
			}
			return nil, fmt.Errorf("failed to find product type: %w", err)
		}
		productTypes[item.ProductTypeUUID] = product
	}
	return productTypes, nil
}
