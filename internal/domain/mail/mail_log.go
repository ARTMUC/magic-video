package mail

import "github.com/ARTMUC/magic-video/internal/domain/base"

type MailLog struct {
	base.BaseModel

	RecipientName  string `gorm:"type:varchar(255)"`
	RecipientEmail string `gorm:"type:varchar(255)"`
	Template       string `gorm:"type:varchar(255)"`
	Status         string `gorm:"type:varchar(255)"`
	Error          string `gorm:"type:text"`
	Reference      string `gorm:"type:varchar(255)"`
}
