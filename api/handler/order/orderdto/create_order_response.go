package dto

type CreateOrderResponse struct {
	Body CreateOrderRequestBody
}

type CreateOrderResponseBody struct {
	Success    bool   `json:"success" example:"true"`
	SessionID  string `json:"sessionId" example:"order_1640995200"`
	Token      string `json:"token" example:"ABC123DEF456"`
	PaymentURL string `json:"paymentUrl" example:"https://sandbox.przelewy24.pl/trnRequest/ABC123DEF456"`
}
