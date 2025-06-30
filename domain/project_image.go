package domain

import "time"

type ProjectImage struct {
	ID        string
	IDProject string
	FileName  string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProjectImagePayload struct {
	IDProject string `json:"id_project" form:"id_project" validate:"required"`
	FileName  string `json:"file_name" form:"file_name" validate:"required"`
	Url       string `json:"url" form:"url" validate:"required,url"`
}

type UpdateProjectImagePayload struct {
	FileName *string `json:"file_name,omitempty" form:"file_name" validate:"omitempty"`
	Url      *string `json:"url,omitempty" form:"url" validate:"omitempty,url"`
}
