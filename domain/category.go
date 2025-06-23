package domain

import "time"

type Category struct {
	Id          string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
