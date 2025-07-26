package mailer

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/mail"
)

type CustomerAccessEmailSender interface {
	Send(to EmailRecipient, apiLink string) (*mail.MailLog, error)
}

type customerAccessEmailSender struct {
	config            config.ServerConfig
	sender            EmailSender
	mailLogRepository mail.MailLogRepository
}

func NewCustomerAccessEmailSender(
	config config.ServerConfig,
	sender EmailSender,
	mailLogRepository mail.MailLogRepository,
) CustomerAccessEmailSender {
	return &customerAccessEmailSender{
		config:            config,
		sender:            sender,
		mailLogRepository: mailLogRepository,
	}
}

func (s *customerAccessEmailSender) Send(to EmailRecipient, apiLink string) (*mail.MailLog, error) {
	template := `<!DOCTYPE html>
<html lang="pl">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Dostęp do Magiczny Prezent</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f6f8;
      margin: 0;
      padding: 0;
      color: #333;
    }
    .container {
      max-width: 600px;
      margin: 30px auto;
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
      padding: 30px;
    }
    h1 {
      color: #2c3e50;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      font-size: 16px;
      line-height: 1.5;
      margin-bottom: 25px;
    }
    a.button {
      display: inline-block;
      padding: 12px 25px;
      background-color: #3498db;
      color: white !important;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    a.button:hover {
      background-color: #2980b9;
    }
    .footer {
      margin-top: 30px;
      font-size: 12px;
      color: #777;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Witamy, {{customer_name}}!</h1>
    <p>Miło nam udostępnić Ci dostęp do naszej platformy. Możesz rozpocząć, klikając poniższy przycisk, aby uzyskać dostęp do API:</p>
    <p><a href="{{api_link}}" class="button" target="_blank" rel="noopener">Dostęp do aplikacji</a></p>
    <p>Jeśli przycisk powyżej nie działa, skopiuj i wklej poniższy link do przeglądarki:</p>
    <p><a href="{{api_link}}" target="_blank" rel="noopener" style="color:#3498db;">{{api_link}}</a></p>
    <p>W razie pytań lub potrzeby pomocy, prosimy o kontakt.</p>
    <p>Pozdrawiamy,<br>Zespół Magiczny Prezent</p>
    <div class="footer">
      © {{current_year}} Magiczny Prezent. Wszelkie prawa zastrzeżone.
    </div>
  </div>
</body>
</html>
`
	template = strings.ReplaceAll(template, "{{customer_name}}", to.Name)
	template = strings.ReplaceAll(template, "{{current_year}}", time.Now().Format("2006"))
	template = strings.ReplaceAll(template, "{{api_link}}", apiLink)

	mailLog, err := s.mailLogRepository.FindOne(base.ReadOptions{
		Scopes: []base.Scope{
			mail.MailLogScopes{}.WithEmail(to.Email),
			mail.MailLogScopes{}.WithTemplate("CustomerAccess"),
			mail.MailLogScopes{}.OrderBy("id DESC"),
		},
	})
	if err != nil {
		if errors.Is(err, base.ErrRecordNotFound) {
			//ok
		} else {
			return nil, fmt.Errorf("failed to list mail logs: %w", err)
		}
	} else {
		return mailLog, nil
	}

	mailLog, err = s.sender.Send(EmailRequest{
		To:           EmailRecipient{to.Name, to.Email},
		Subject:      "Dostęp do Magiczny Prezent",
		HTMLContent:  template,
		TemplateName: "CustomerAccess",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return mailLog, nil
}
