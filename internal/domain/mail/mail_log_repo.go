package mail

import (
	"github.com/ARTMUC/magic-video/internal/domain/base"
	"gorm.io/gorm"
)

type MailLogRepository interface {
	base.BaseRepository[MailLog]
}

type mailLogRepo struct {
	*base.BaseRepo[MailLog]

	db *gorm.DB
}

func NewMailLogRepo(db *gorm.DB) MailLogRepository {
	return &mailLogRepo{base.NewBaseRepository[MailLog](db), db}
}

type MailLogScopes struct {
	base.BaseScopes
}

func (u MailLogScopes) WithEmail(email string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("mail_logs.email = ?", email)
	}
}

func (u MailLogScopes) WithTemplate(template string) base.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("mail_logs.template = ?", template)
	}
}
