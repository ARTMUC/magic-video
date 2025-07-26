package p24

// Struktura do rejestracji transakcji
type TransactionRegisterRequest struct {
	MerchantId       int         `json:"merchantId"`
	PosId            int         `json:"posId"`
	SessionId        string      `json:"sessionId"`
	Amount           int         `json:"amount"`
	Currency         string      `json:"currency"`
	Description      string      `json:"description"`
	Email            string      `json:"email"`
	Client           string      `json:"client,omitempty"`
	Address          string      `json:"address,omitempty"`
	Zip              string      `json:"zip,omitempty"`
	City             string      `json:"city,omitempty"`
	Country          string      `json:"country"`
	Phone            string      `json:"phone,omitempty"`
	Language         string      `json:"language,omitempty"`
	Method           int         `json:"method,omitempty"`
	UrlReturn        string      `json:"urlReturn"`
	UrlStatus        string      `json:"urlStatus"`
	TimeLimit        int         `json:"timeLimit,omitempty"`
	Channel          int         `json:"channel,omitempty"`
	WaitForResult    bool        `json:"waitForResult,omitempty"`
	RegulationAccept bool        `json:"regulationAccept"`
	Shipping         int         `json:"shipping,omitempty"`
	Sign             string      `json:"sign"`
	Encoding         string      `json:"encoding,omitempty"`
	MethodRefId      string      `json:"methodRefId,omitempty"`
	Cart             []CartItem  `json:"cart,omitempty"`
	Additional       *Additional `json:"additional,omitempty"`
}

// Struktura elementu koszyka
type CartItem struct {
	SellerId       string `json:"sellerId,omitempty"`
	SellerCategory string `json:"sellerCategory,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Quantity       int    `json:"quantity"`
	Price          int    `json:"price"`
	Number         string `json:"number,omitempty"`
}

// Dodatkowe informacje do transakcji
type Additional struct {
	Shipping *Shipping `json:"shipping,omitempty"`
}

// Dane wysyłki
type Shipping struct {
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
	Zip     string `json:"zip,omitempty"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
}

// Odpowiedź rejestracji transakcji
type TransactionRegisterResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	ResponseCode int    `json:"responseCode"`
	Error        string `json:"error,omitempty"`
}

// Struktura do weryfikacji transakcji
type TransactionVerifyRequest struct {
	MerchantId int    `json:"merchantId"`
	PosId      int    `json:"posId"`
	SessionId  string `json:"sessionId"`
	Amount     int    `json:"amount"`
	Currency   string `json:"currency"`
	OrderId    int    `json:"orderId"`
	Sign       string `json:"sign"`
}

// Odpowiedź weryfikacji transakcji
type TransactionVerifyResponse struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
	ResponseCode int    `json:"responseCode"`
	Error        string `json:"error,omitempty"`
}

// Struktura powiadomienia (notyfikacji)
type NotificationRequest struct {
	MerchantId   int    `json:"merchantId"`
	PosId        int    `json:"posId"`
	SessionId    string `json:"sessionId"`
	Amount       int    `json:"amount"`
	OriginAmount int    `json:"originAmount"`
	Currency     string `json:"currency"`
	OrderId      int    `json:"orderId"`
	MethodId     int    `json:"methodId"`
	Statement    string `json:"statement"`
	Sign         string `json:"sign"`
}
