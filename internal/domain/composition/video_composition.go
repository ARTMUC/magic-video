package domain

const (
	VideoCompositionStatusCompleted  = "completed"
	VideoCompositionStatusInProgress = "in_progress"
	VideoCompositionStatusFailed     = "failed"
	VideoCompositionStatusCancelled  = "cancelled"
	VideoCompositionStatusCreated    = "created"
)

type VideoComposition struct {
	BaseModel

	CustomerID    uint
	Customer      *Customer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VideoTemplate string    `gorm:"type:varchar(255)"` // for example: superhero, pirates, santaclaus ...
	Status        string    `gorm:"type:varchar(255)"`
	Images        []Image
}
