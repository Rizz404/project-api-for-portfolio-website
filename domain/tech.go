package domain

import "time"

type Tech struct {
	Id          string
	Name        string
	Description *string
	LogoURL     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
