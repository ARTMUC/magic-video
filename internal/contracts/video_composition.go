package contracts

import "github.com/google/uuid"

type VideoComposition struct {
	BaseModel

	CustomerID    uuid.UUID
	VideoTemplate string
	Status        string
	Images        []Image
}

type Image struct {
	BaseModel
	
	Name            string
	PresetImageType string
}
