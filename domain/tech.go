package domain

import "time"

type Tech struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	LogoURL     *string   `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTechPayload struct {
	Name        string  `json:"name" form:"name" validate:"required,min=3,max=50"`
	Description *string `json:"description,omitempty" form:"description" validate:"omitempty,min=3,max=500"`
	LogoURL     *string `json:"logo_url,omitempty" form:"logo_url" validate:"omitempty,url"`
}

type UpdateTechPayload struct {
	Name        *string `json:"name,omitempty" form:"name" validate:"omitempty,min=3,max=50"`
	Description *string `json:"description,omitempty" form:"description" validate:"omitempty,min=3,max=500"`
	LogoURL     *string `json:"logo_url,omitempty" form:"logo_url" validate:"omitempty,url"`
}
