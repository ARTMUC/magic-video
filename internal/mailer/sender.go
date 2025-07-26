package mailer

import (
	"github.com/ARTMUC/magic-video/internal/domain/mail"
)

type EmailSender interface {
	Send(EmailRequest) (*mail.MailLog, error)
}

type EmailRecipient struct {
	Name  string
	Email string
}
type EmailRequest struct {
	To           EmailRecipient
	Subject      string
	HTMLContent  string
	TemplateName string
}
