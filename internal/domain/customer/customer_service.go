package service

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain"
	"github.com/ARTMUC/magic-video/internal/mailer"
	"github.com/ARTMUC/magic-video/internal/pkg/crypto"
	"github.com/ARTMUC/magic-video/internal/repository"
	"github.com/google/uuid"
)

var ErrCustomerAccessNotFound = errors.New("customer access not found")
var ErrCustomerNotFound = errors.New("customer not found")

type CustomerService interface {
	CreateAccessThruEmail(email string) (*domain.Customer, error)
	GetCustomerFromToken(customerUUID uuid.UUID, accessToken string) (*domain.Customer, error)
	GetCustomerByUUID(customerUUID uuid.UUID) (*domain.Customer, error)
}

type customerService struct {
	customerRepository       repository.CustomerRepository
	customerAccessRepository repository.CustomerAccessRepository
	customerAccessMailSender mailer.CustomerAccessEmailSender
	serverConfig             config.ServerConfig
	encryptionConfig         config.EncryptionConfig
}

func NewCustomerService(
	customerRepository repository.CustomerRepository,
	customerAccessRepository repository.CustomerAccessRepository,
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

func (s *customerService) GetCustomerByUUID(customerUUID uuid.UUID) (*domain.Customer, error) {
	customer, err := s.customerRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{repository.WithUUID(customerUUID)},
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to find customer in db: %w", err)
	}
	return customer, nil
}

func (s *customerService) CreateAccessThruEmail(email string) (*domain.Customer, error) {
	customer, err := s.customerRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{
			repository.CustomerScopes{}.WithEmail(email),
		},
	})
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find customer with email %s: %w", email, err)
		} else {
			customer = &domain.Customer{
				Email: email,
			}
			err = s.customerRepository.Create(repository.WriteOptions{}, customer)
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

func (s *customerService) CreateCustomerAccess(customer *domain.Customer) (*domain.CustomerAccess, error) {
	customerAccess, err := s.customerAccessRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{
			repository.CustomerAccessScopes{}.WithCustomer(customer),
			repository.CustomerAccessScopes{}.WithNotExpired(),
			repository.CustomerAccessScopes{}.OrderBy("customer_accesses.id DESC"),
		},
	})
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
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

	customerAccess = &domain.CustomerAccess{
		CustomerID:      customer.ID,
		Customer:        customer,
		AccessToken:     encodedToken,
		TokenExpireDate: expiresAt,
	}

	err = s.customerAccessRepository.Create(repository.WriteOptions{}, customerAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to save customer access in db: %w", err)
	}

	return customerAccess, nil
}

func (s *customerService) GetCustomerFromToken(customerUUID uuid.UUID, token string) (*domain.Customer, error) {
	customer, err := s.customerRepository.FindOne(repository.ReadOptions{
		Scopes: []repository.Scope{repository.WithUUID(customerUUID)},
	})
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find customer in db: %w", err)
		}
		return nil, ErrCustomerNotFound
	}

	customerAccesses, err := s.customerAccessRepository.FindMany(repository.ReadOptions{
		Scopes: []repository.Scope{
			repository.CustomerAccessScopes{}.WithNotExpired(),
			repository.CustomerAccessScopes{}.WithCustomer(customer),
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
