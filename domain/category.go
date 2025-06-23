package domain

import "time"

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryPayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateCategoryPayload struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
