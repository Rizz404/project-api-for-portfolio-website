package domain

import "time"

type Language struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LangCode  string    `json:"lang_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateLanguagePayload struct {
	Name     string `json:"name" form:"name" validate:"required,min=3,max=30,alphanum"`
	LangCode string `json:"lang_code" form:"lang_code" validate:"required,min=3,max=30,alphanum"`
}

type UpdateLanguagePayload struct {
	Name     *string `json:"name,omitempty" form:"name" validate:"omitempty,min=3,max=30,alphanum"`
	LangCode *string `json:"lang_code,omitempty" form:"lang_code" validate:"omitempty,min=3,max=30,alphanum"`
}
