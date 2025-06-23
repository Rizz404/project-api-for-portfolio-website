package domain

import "time"

type ProjectImage struct {
	ID        string
	IDProject string
	FileName  string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
