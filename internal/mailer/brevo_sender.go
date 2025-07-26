package mailer

import (
	"context"
	"fmt"

	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"github.com/ARTMUC/magic-video/internal/domain/mail"
	brevo "github.com/getbrevo/brevo-go/lib"
)

type brevoEmailSender struct {
	config            config.BrevoEmailClientConfig
	mailLogRepository mail.MailLogRepository
}

func NewBrevoEmailSender(
	config config.BrevoEmailClientConfig,
	mailLogRepository mail.MailLogRepository,
) EmailSender {
	return &brevoEmailSender{
		config:            config,
		mailLogRepository: mailLogRepository,
	}
}

func (s *brevoEmailSender) Send(input EmailRequest) (*mail.MailLog, error) {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", s.config.ApiKey())
	cfg.AddDefaultHeader("partner-key", s.config.ApiKey())
	client := brevo.NewAPIClient(cfg)

	sender := brevo.SendSmtpEmailSender{
		Name:  s.config.SenderName(),
		Email: s.config.SenderEmail(),
	}

	email := brevo.SendSmtpEmail{
		Sender: &sender,
		To: []brevo.SendSmtpEmailTo{
			{
				Email: input.To.Email,
				Name:  input.To.Name,
			},
		},
		Subject:     input.Subject,
		HtmlContent: input.HTMLContent,
	}

	mailLog := &mail.MailLog{
		RecipientName:  input.To.Name,
		RecipientEmail: input.To.Email,
		Template:       input.TemplateName,
	}

	resp, _, err := client.TransactionalEmailsApi.SendTransacEmail(context.Background(), email)
	if err != nil {
		if x, ok := err.(interface {
			Model() interface{}
		}); ok && x.Model() != nil {
			err = fmt.Errorf("error sending email: %w, resp: %s", err, x.Model())
		} else {
			err = fmt.Errorf("error sending email: %w", err)
		}
		mailLog.Status = "error"
		mailLog.Error = err.Error()
		mailLogErr := s.mailLogRepository.Create(base.WriteOptions{}, mailLog)
		if mailLogErr != nil {
			return nil, fmt.Errorf("failed to save mail log in db: %w: %w", mailLogErr, err)
		}
		return nil, err
	}

	mailLog.Status = "success"
	mailLog.Reference = resp.MessageId
	mailLogErr := s.mailLogRepository.Create(base.WriteOptions{}, mailLog)
	if mailLogErr != nil {
		return nil, fmt.Errorf("failed to save mail log in db: %w: %w", mailLogErr, err)
	}

	return mailLog, nil
}
