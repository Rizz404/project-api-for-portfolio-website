package domain

import "time"

type Tech struct {
	ID          string
	Name        string
	Description *string
	LogoURL     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
