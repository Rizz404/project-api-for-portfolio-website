package domain

import "time"

type TechStack struct {
	IDProject string
	IDTech    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
