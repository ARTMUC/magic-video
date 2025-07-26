package domain

type Image struct {
	BaseModel

	Name               string `gorm:"type:varchar(255)"`
	PresetImageType    string `gorm:"type:varchar(255)"` // for example front, side || kid, parent || man, woman || person1, person2, person3 - whatever just to recognize it in video processor
	VideoCompositionID uint
	VideoComposition   *VideoComposition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT""`
}
