package domain

import "time"

type ProjectTranslation struct {
	ID          string    `json:"id"`
	IDProject   string    `json:"id_project"`
	IDLanguage  string    `json:"id_language"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreateProjectTranslationPayload struct {
	IDProject   string  `json:"id_project" form:"id_project" validate:"required,ulid"`
	IDLanguage  string  `json:"id_language" form:"id_language" validate:"required,ulid"`
	Name        string  `json:"name" form:"name" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" form:"description" validate:"omitempty,min=3"`
}

type UpdateProjectTranslationPayload struct {
	IDLanguage  string  `json:"id_language,omitempty" form:"id_language" validate:"required,ulid"`
	Name        string  `json:"name,omitempty" form:"name" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" form:"description" validate:"omitempty,min=3"`
}
