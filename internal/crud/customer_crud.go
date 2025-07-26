package crud

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"github.com/ARTMUC/magic-video/internal/repository"
)

type MailLogCrud struct {
	BaseCrud[domain.MailLog]
	repository repository.MailLogRepository
}

func NewMailLogCrud(
	repository repository.MailLogRepository,
) *MailLogCrud {
	return &MailLogCrud{BaseCrud: newBaseCrud(repository), repository: repository}
}
