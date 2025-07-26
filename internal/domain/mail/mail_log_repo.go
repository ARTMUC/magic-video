package repository

import (
	"github.com/ARTMUC/magic-video/internal/domain"
	"gorm.io/gorm"
)

type MailLogRepository interface {
	BaseRepository[domain.MailLog]
}

type mailLogRepo struct {
	*BaseRepo[domain.MailLog]

	db *gorm.DB
}

func NewMailLogRepo(db *gorm.DB) MailLogRepository {
	return &mailLogRepo{NewBaseRepository[domain.MailLog](db), db}
}

type MailLogScopes struct {
	BaseScopes
}

func (u MailLogScopes) WithEmail(email string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("mail_logs.email = ?", email)
	}
}

func (u MailLogScopes) WithTemplate(template string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("mail_logs.template = ?", template)
	}
}
