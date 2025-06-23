package domain

import "time"

type ProjectImage struct {
	Id        string
	IdProject string
	FileName  string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
