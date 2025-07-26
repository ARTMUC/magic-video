package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ARTMUC/magic-video/internal/config"
)

var apiKey = "your_brevo_api_key"
var url = "https://api.brevo.com/v3/smtp/email"

type EmailSender interface {
	Send(input EmailRequest) error
}

type emailSender struct {
	config config.BrevoEmailClientConfig
}

func NewEmailSender(config config.BrevoEmailClientConfig) EmailSender {
	return &emailSender{config: config}
}

type EmailRequest struct {
	Sender      map[string]string   `json:"sender"`
	To          []map[string]string `json:"to"`
	Subject     string              `json:"subject"`
	HTMLContent string              `json:"htmlContent"`
}

func (e *emailSender) Send(input EmailRequest) error {
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}
