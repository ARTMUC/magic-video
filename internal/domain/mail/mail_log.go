package domain

type MailLog struct {
	BaseModel

	RecipientName  string `gorm:"type:varchar(255)"`
	RecipientEmail string `gorm:"type:varchar(255)"`
	Template       string `gorm:"type:varchar(255)"`
	Status         string `gorm:"type:varchar(255)"`
	Error          string `gorm:"type:text"`
	Reference      string `gorm:"type:varchar(255)"`
}
