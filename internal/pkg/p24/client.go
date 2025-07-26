package p24

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Environment constants
const (
	EnvLive    = "live"
	EnvSandbox = "sandbox"
)

// API endpoints
const (
	LiveURL    = "https://secure.przelewy24.pl"
	SandboxURL = "https://sandbox.przelewy24.pl"
)

// Client represents the Przelewy24 API client
type Client struct {
	MerchantID  string
	PosID       string
	Salt        string
	APIKey      string
	Environment string
	BaseURL     string
	HTTPClient  *http.Client
}

// NewClient creates a new Przelewy24 client
func NewClient(merchantID, posID string, salt, apiKey, environment string) *Client {
	baseURL := LiveURL
	if environment == EnvSandbox {
		baseURL = SandboxURL
	}

	return &Client{
		MerchantID:  merchantID,
		PosID:       posID,
		Salt:        salt,
		APIKey:      apiKey,
		Environment: environment,
		BaseURL:     baseURL,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// TransactionRequest represents a transaction registration request
type TransactionRequest struct {
	SessionID     string `json:"p24_session_id"`
	Amount        int    `json:"p24_amount"`
	Currency      string `json:"p24_currency"`
	Description   string `json:"p24_description"`
	Email         string `json:"p24_email"`
	Country       string `json:"p24_country"`
	URLReturn     string `json:"p24_url_return"`
	URLStatus     string `json:"p24_url_status"`
	TimeLimit     int    `json:"p24_time_limit,omitempty"`
	Channel       int    `json:"p24_channel,omitempty"`
	Method        int    `json:"p24_method,omitempty"`
	Language      string `json:"p24_language,omitempty"`
	ClientName    string `json:"p24_client,omitempty"`
	Address       string `json:"p24_address,omitempty"`
	Zip           string `json:"p24_zip,omitempty"`
	City          string `json:"p24_city,omitempty"`
	Phone         string `json:"p24_phone,omitempty"`
	TransferLabel string `json:"p24_transfer_label,omitempty"`
	Shipping      int    `json:"p24_shipping,omitempty"`
	Encoding      string `json:"p24_encoding,omitempty"`
}

// TransactionResponse represents the response from transaction registration
type TransactionResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	ReturnErrorCode    int    `json:"returnErrorCode"`
	ReturnErrorMessage string `json:"returnErrorMessage"`
}

// VerificationRequest represents a transaction verification request
type VerificationRequest struct {
	SessionID string `json:"p24_session_id"`
	Amount    int    `json:"p24_amount"`
	Currency  string `json:"p24_currency"`
	OrderID   int    `json:"p24_order_id"`
}

// VerificationResponse represents the response from transaction verification
type VerificationResponse struct {
	Error string `json:"error"`
}

// TestConnectionResponse represents the response from test connection
type TestConnectionResponse struct {
	Error string `json:"error"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ImgURL string `json:"imgUrl"`
	Mobile bool   `json:"mobile"`
}

// PaymentMethodsResponse represents the response from payment methods endpoint
type PaymentMethodsResponse struct {
	Data  []PaymentMethod `json:"data"`
	Error string          `json:"error"`
}

// RegisterTransaction registers a new transaction
func (c *Client) RegisterTransaction(req TransactionRequest) (*TransactionResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)
	formData.Set("p24_session_id", req.SessionID)
	formData.Set("p24_amount", strconv.Itoa(req.Amount))
	formData.Set("p24_currency", req.Currency)
	formData.Set("p24_description", req.Description)
	formData.Set("p24_email", req.Email)
	formData.Set("p24_country", req.Country)
	formData.Set("p24_url_return", req.URLReturn)
	formData.Set("p24_url_status", req.URLStatus)
	formData.Set("p24_api_version", "3.2")

	// Optional fields
	if req.TimeLimit > 0 {
		formData.Set("p24_time_limit", strconv.Itoa(req.TimeLimit))
	}
	if req.Channel > 0 {
		formData.Set("p24_channel", strconv.Itoa(req.Channel))
	}
	if req.Method > 0 {
		formData.Set("p24_method", strconv.Itoa(req.Method))
	}
	if req.Language != "" {
		formData.Set("p24_language", req.Language)
	}
	if req.ClientName != "" {
		formData.Set("p24_client", req.ClientName)
	}
	if req.Address != "" {
		formData.Set("p24_address", req.Address)
	}
	if req.Zip != "" {
		formData.Set("p24_zip", req.Zip)
	}
	if req.City != "" {
		formData.Set("p24_city", req.City)
	}
	if req.Phone != "" {
		formData.Set("p24_phone", req.Phone)
	}
	if req.TransferLabel != "" {
		formData.Set("p24_transfer_label", req.TransferLabel)
	}
	if req.Shipping > 0 {
		formData.Set("p24_shipping", strconv.Itoa(req.Shipping))
	}
	if req.Encoding != "" {
		formData.Set("p24_encoding", req.Encoding)
	}

	// Generate signature
	signData := fmt.Sprintf("%s|%s|%d|%s|%s|%s|%s",
		req.SessionID, c.PosID, req.Amount, req.Currency, req.Description, c.MerchantID, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request
	result, _, err := postForm[TransactionResponse](c.HTTPClient, c.BaseURL, "/api/v1/transaction/register", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to post form: %w", err)
	}

	return &result, nil
}

// VerifyTransaction verifies a completed transaction
func (c *Client) VerifyTransaction(req VerificationRequest) (*VerificationResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)
	formData.Set("p24_session_id", req.SessionID)
	formData.Set("p24_amount", strconv.Itoa(req.Amount))
	formData.Set("p24_currency", req.Currency)
	formData.Set("p24_order_id", strconv.Itoa(req.OrderID))

	// Generate signature
	signData := fmt.Sprintf("%s|%d|%d|%s|%s",
		req.SessionID, req.OrderID, req.Amount, req.Currency, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request
	result, _, err := postForm[VerificationResponse](c.HTTPClient, c.BaseURL, "/api/v1/transaction/verify", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to post form: %w", err)
	}

	return &result, nil
}

// TestConnection tests the connection to Przelewy24 API
func (c *Client) TestConnection() (*TestConnectionResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)

	// Generate signature
	signData := fmt.Sprintf("%s|%s|%s", c.PosID, c.MerchantID, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request
	result, _, err := postForm[TestConnectionResponse](c.HTTPClient, c.BaseURL, "/api/v1/testConnection", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to post form: %w", err)
	}

	return &result, nil
}

// GetPaymentMethods retrieves available payment methods
func (c *Client) GetPaymentMethods(lang string) (*PaymentMethodsResponse, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)
	if lang != "" {
		formData.Set("p24_lang", lang)
	}

	// Generate signature
	signData := fmt.Sprintf("%s|%s|%s", c.PosID, c.MerchantID, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request
	result, _, err := postForm[PaymentMethodsResponse](c.HTTPClient, c.BaseURL, "/api/v1/payment/methods", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to post form: %w", err)
	}

	return &result, nil
}

// GetTransactionBySessionID retrieves transaction details by session ID
func (c *Client) GetTransactionBySessionID(sessionID string) (map[string]interface{}, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)
	formData.Set("p24_session_id", sessionID)

	// Generate signature
	signData := fmt.Sprintf("%s|%s|%s|%s", sessionID, c.PosID, c.MerchantID, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request
	result, _, err := postForm[map[string]any](c.HTTPClient, c.BaseURL, "/api/v1/transaction/by/session", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to post form: %w", err)
	}

	return result, nil
}

// RefundTransaction creates a refund for a transaction
func (c *Client) RefundTransaction(orderID int, sessionID string, amount int, description string) (map[string]interface{}, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("p24_merchant_id", c.MerchantID)
	formData.Set("p24_pos_id", c.PosID)
	formData.Set("p24_order_id", strconv.Itoa(orderID))
	formData.Set("p24_session_id", sessionID)
	formData.Set("p24_amount", strconv.Itoa(amount))
	formData.Set("p24_description", description)

	// Generate signature for refund
	signData := fmt.Sprintf("%s|%s|%s|%s|%s", sessionID, orderID, amount, c.MerchantID, c.Salt)
	formData.Set("p24_sign", c.md5Hash(signData))

	// Make request using API key authentication
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/transaction/refund", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.PosID, c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetPaymentURL generates the payment URL for a registered transaction
func (c *Client) GetPaymentURL(token string) string {
	return fmt.Sprintf("%s/trnRequest/%s", c.BaseURL, token)
}

// VerifyWebhookSignature verifies the signature of incoming webhook data
func (c *Client) VerifyWebhookSignature(data *NotificationRequest) bool {
	checksumData := map[string]interface{}{
		"sessionId": data.SessionId,
		"orderId":   data.OrderId,
		"amount":    data.Amount,
		"currency":  data.Currency,
		"crc":       c.Salt,
	}

	jsonBytes, err := json.Marshal(checksumData)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return false
	}

	hash := sha512.Sum384(jsonBytes)
	expectedSign := hex.EncodeToString(hash[:])
	
	return expectedSign == data.Sign
}

// Helper methods

// postForm sends a POST request with form data
func postForm[R any](httpClient *http.Client, baseUrl string, endpoint string, data url.Values) (R, *http.Response, error) {
	var r R

	req, err := http.NewRequest("POST", baseUrl+endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return r, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return r, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return r, nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return r, resp, nil
}

// md5Hash generates MD5 hash of the input string
func (c *Client) md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Utility functions for common operations

// ParseWebhookData parses webhook data from HTTP request
func ParseWebhookData(r *http.Request) (*NotificationRequest, error) {
	defer r.Body.Close()
	var data NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse webhook JSON: %w", err)
	}

	return &data, nil
}

// FormatFloatAmount converts amount from float to integer (multiply by 100 for grosze)
func FormatFloatAmount(amount float64) int {
	return int(amount * 100)
}

// FormatFloatAmount converts amount from float to integer (multiply by 100 for grosze)
func FormatDecimalAmount(amount decimal.Decimal) int {
	return int(amount.Mul(decimal.NewFromInt(100)).IntPart())
}

// ParseAmount converts integer amount back to float (divide by 100)
func ParseAmount(amount int) float64 {
	return float64(amount) / 100
}

// Example usage and integration helpers

// TransactionBuilder helps build transaction requests
type TransactionBuilder struct {
	req TransactionRequest
}

// NewTransactionBuilder creates a new transaction builder
func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{
		req: TransactionRequest{
			Currency: "PLN",
			Country:  "PL",
			Language: "pl",
			Encoding: "UTF-8",
		},
	}
}

// SetSessionID sets the session ID
func (tb *TransactionBuilder) SetSessionID(sessionID string) *TransactionBuilder {
	tb.req.SessionID = sessionID
	return tb
}

// SetAmount sets the amount in grosze
func (tb *TransactionBuilder) SetAmount(amount int) *TransactionBuilder {
	tb.req.Amount = amount
	return tb
}

// SetAmountFloat sets the amount from float (PLN)
func (tb *TransactionBuilder) SetAmountFloat(amount float64) *TransactionBuilder {
	tb.req.Amount = FormatFloatAmount(amount)
	return tb
}

// SetAmountFloat sets the amount from float (PLN)
func (tb *TransactionBuilder) SetAmountDecimal(amount decimal.Decimal) *TransactionBuilder {
	tb.req.Amount = FormatDecimalAmount(amount)
	return tb
}

// SetDescription sets the transaction description
func (tb *TransactionBuilder) SetDescription(description string) *TransactionBuilder {
	tb.req.Description = description
	return tb
}

// SetEmail sets the customer email
func (tb *TransactionBuilder) SetEmail(email string) *TransactionBuilder {
	tb.req.Email = email
	return tb
}

// SetReturnURL sets the return URL
func (tb *TransactionBuilder) SetReturnURL(url string) *TransactionBuilder {
	tb.req.URLReturn = url
	return tb
}

// SetStatusURL sets the status URL for webhooks
func (tb *TransactionBuilder) SetStatusURL(url string) *TransactionBuilder {
	tb.req.URLStatus = url
	return tb
}

// SetCustomerDetails sets customer details
func (tb *TransactionBuilder) SetCustomerDetails(name, address, zip, city, phone string) *TransactionBuilder {
	tb.req.ClientName = name
	tb.req.Address = address
	tb.req.Zip = zip
	tb.req.City = city
	tb.req.Phone = phone
	return tb
}

// Build returns the built transaction request
func (tb *TransactionBuilder) Build() TransactionRequest {
	return tb.req
}
