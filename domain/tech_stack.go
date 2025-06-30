package domain

import "time"

type TechStack struct {
	IDProject string
	IDTech    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTechStackPayload struct {
	IDProject string `json:"id_project" form:"id_project" validate:"required,ulid"`
	IDTech    string `json:"id_tech" form:"id_tech" validate:"required,ulid"`
}

type UpdateTechStackPayload struct {
	IDProject *string `json:"id_project,omitempty" form:"id_project" validate:"ulid,omitempty"`
	IDTech    *string `json:"id_tech,omitempty" form:"id_tech" validate:"ulid,omitempty"`
}

type DeleteTechStackPayload struct {
	IDProject string `json:"id_project" form:"id_project" validate:"required,ulid"`
	IDTech    string `json:"id_tech" form:"id_tech" validate:"required,ulid"`
}
