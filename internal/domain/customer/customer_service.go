package customer

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/mailer"
	"github.com/ARTMUC/magic-video/internal/pkg/crypto"
	"github.com/google/uuid"
)

var ErrCustomerAccessNotFound = errors.New("customer access not found")
var ErrCustomerNotFound = errors.New("customer not found")

type CustomerService interface {
	CreateAccessThruEmail(email string) (*Customer, error)
	GetCustomerFromToken(customerUUID uuid.UUID, accessToken string) (*Customer, error)
	GetCustomerByUUID(customerUUID uuid.UUID) (*Customer, error)
}

type customerService struct {
	customerRepository       CustomerRepository
	customerAccessRepository CustomerAccessRepository
	customerAccessMailSender mailer.CustomerAccessEmailSender
	serverConfig             config.ServerConfig
	encryptionConfig         config.EncryptionConfig
}

func NewCustomerService(
	customerRepository CustomerRepository,
	customerAccessRepository CustomerAccessRepository,
	customerAccessMailSender mailer.CustomerAccessEmailSender,
	serverConfig config.ServerConfig,
	encryptionConfig config.EncryptionConfig,
) CustomerService {
	return &customerService{
		customerRepository:       customerRepository,
		customerAccessRepository: customerAccessRepository,
		customerAccessMailSender: customerAccessMailSender,
		serverConfig:             serverConfig,
		encryptionConfig:         encryptionConfig,
	}
}

func (s *customerService) GetCustomerByUUID(customerUUID uuid.UUID) (*Customer, error) {
	customer, err := s.customerRepository.FindOne(base.ReadOptions{
		Scopes: []base.Scope{base.WithUUID(customerUUID)},
	})
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to find customer in db: %w", err)
	}
	return customer, nil
}

func (s *customerService) CreateAccessThruEmail(email string) (*Customer, error) {
	customer, err := s.customerRepository.FindOne(base.ReadOptions{
		Scopes: []base.Scope{
			CustomerScopes{}.WithEmail(email),
		},
	})
	if err != nil {
		if !errors.Is(err, base.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find customer with email %s: %w", email, err)
		} else {
			customer = &Customer{
				Email: email,
			}
			err = s.customerRepository.Create(base.WriteOptions{}, customer)
			if err != nil {
				return nil, fmt.Errorf("failed to create customer in db: %w", err)
			}
		}
	}

	customerAccess, err := s.CreateCustomerAccess(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer access: %w", err)
	}

	recipient := mailer.EmailRecipient{
		Name:  "Nowy kliencie",
		Email: customer.Email,
	}
	if customer.Name.Valid {
		recipient.Name = customer.Name.V
	}

	token, err := crypto.DecryptToken(customerAccess.AccessToken, s.encryptionConfig.EncryptionKey())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	_, err = s.customerAccessMailSender.Send(
		recipient,
		path.Join(s.serverConfig.WebsiteUrl(), token),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send customer access email: %w", err)
	}

	return customer, nil
}

func (s *customerService) CreateCustomerAccess(customer *Customer) (*CustomerAccess, error) {
	customerAccess, err := s.customerAccessRepository.FindOne(base.ReadOptions{
		Scopes: []base.Scope{
			CustomerAccessScopes{}.WithCustomer(customer),
			CustomerAccessScopes{}.WithNotExpired(),
			CustomerAccessScopes{}.OrderBy("customer_accesses.id DESC"),
		},
	})
	if err != nil {
		if !errors.Is(err, base.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find customer access in db: %w", err)
		} else {
			// ok
		}
	} else {
		return customerAccess, nil
	}
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("token generation failed: %w", err)
	}
	token := strings.ToLower(base32.StdEncoding.EncodeToString(b))
	token = strings.TrimRight(token, "=")

	encodedToken, err := crypto.EncryptToken(token, s.encryptionConfig.EncryptionKey())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	customerAccess = &CustomerAccess{
		CustomerID:      customer.ID,
		Customer:        customer,
		AccessToken:     encodedToken,
		TokenExpireDate: expiresAt,
	}

	err = s.customerAccessRepository.Create(base.WriteOptions{}, customerAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to save customer access in db: %w", err)
	}

	return customerAccess, nil
}

func (s *customerService) GetCustomerFromToken(customerUUID uuid.UUID, token string) (*Customer, error) {
	customer, err := s.customerRepository.FindOne(base.ReadOptions{
		Scopes: []base.Scope{base.WithUUID(customerUUID)},
	})
	if err != nil {
		if !errors.Is(err, base.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find customer in db: %w", err)
		}
		return nil, ErrCustomerNotFound
	}

	customerAccesses, err := s.customerAccessRepository.FindMany(base.ReadOptions{
		Scopes: []base.Scope{
			CustomerAccessScopes{}.WithNotExpired(),
			CustomerAccessScopes{}.WithCustomer(customer),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find customer access in db: %w", err)
	}

	var found bool
	for _, customerAccess := range customerAccesses {
		decodedToken, err := crypto.DecryptToken(customerAccess.AccessToken, s.encryptionConfig.EncryptionKey())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt access token: %w", err)
		}
		if decodedToken == token {
			found = true
		}
	}
	if !found {
		return nil, ErrCustomerAccessNotFound
	}

	return customer, nil
}
