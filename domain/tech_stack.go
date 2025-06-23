package domain

import "time"

type TechStack struct {
	IdProject string
	IdTech    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
