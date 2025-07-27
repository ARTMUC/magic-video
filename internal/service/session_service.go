package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/customer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateCustomerSession(customer *customer.Customer) (*SessionOutput, error)
	ParseCustomerToken(tokenStr string, isRefresh bool) (*JWTClaimsWithEntity[*customer.Customer], error)
	CustomerClaimsFromContext(ctx context.Context) (*JWTClaimsWithEntity[*customer.Customer], bool)
}
type sessionService struct {
	sessionConfig config.SessionConfig
	customerCrud  customer.CustomerCrud
}

func NewSessionService(sessionConfig config.SessionConfig, customerCrud customer.CustomerCrud) SessionService {
	return &sessionService{sessionConfig: sessionConfig, customerCrud: customerCrud}
}

type SessionOutput struct {
	Token        string
	RefreshToken string
}

type JWTClaimsWithEntity[E any] struct {
	JWTClaims

	Entity E
}

type JWTClaims struct {
	jwt.RegisteredClaims

	EntityName string    `json:"entity_name"`
	EntityID   uuid.UUID `json:"entity_id"`
	Refresh    bool      `json:"refresh"`
}

func (s *sessionService) CreateCustomerSession(customer *customer.Customer) (*SessionOutput, error) {
	now := time.Now()
	accessClaims := JWTClaims{
		EntityName: "customer",
		EntityID:   customer.ID,
		Refresh:    false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.sessionConfig.JwtTokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString([]byte(s.sessionConfig.JwtTokenSecret()))
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshClaims := JWTClaims{
		EntityName: "customer",
		EntityID:   customer.ID,
		Refresh:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.sessionConfig.JwtRefreshTokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString([]byte(s.sessionConfig.JwtRefreshTokenSecret()))
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &SessionOutput{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *sessionService) ParseCustomerToken(tokenStr string, isRefresh bool) (*JWTClaimsWithEntity[*customer.Customer], error) {
	secret := s.sessionConfig.JwtTokenSecret()
	if isRefresh {
		secret = s.sessionConfig.JwtRefreshTokenSecret()
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	cus, err := s.customerCrud.Get(claims.EntityID, base.ReadOptions{})
	if err != nil {
		return nil, fmt.Errorf("can't get customer by uuid: %w", err)
	}

	result := &JWTClaimsWithEntity[*customer.Customer]{
		Entity:    cus,
		JWTClaims: *claims,
	}

	return result, nil
}

func (s *sessionService) ParseToken(tokenStr string, isRefresh bool) (*JWTClaims, error) {
	secret := s.sessionConfig.JwtTokenSecret()
	if isRefresh {
		secret = s.sessionConfig.JwtRefreshTokenSecret()
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *sessionService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.ParseToken(refreshToken, true)
	if err != nil {
		return "", fmt.Errorf("parsing refresh token: %w", err)
	}

	if !claims.Refresh {
		return "", fmt.Errorf("provided token is not a refresh token")
	}

	now := time.Now()
	newAccessClaims := JWTClaims{
		EntityName: claims.EntityName,
		EntityID:   claims.EntityID,
		Refresh:    false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.sessionConfig.JwtTokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	signedToken, err := token.SignedString(s.sessionConfig.JwtTokenSecret())
	if err != nil {
		return "", fmt.Errorf("signing new access token: %w", err)
	}

	return signedToken, nil
}

func (s *sessionService) CustomerClaimsFromContext(ctx context.Context) (*JWTClaimsWithEntity[*customer.Customer], bool) {
	claims, ok := ctx.Value("auth").(*JWTClaimsWithEntity[*customer.Customer])
	return claims, ok
}
