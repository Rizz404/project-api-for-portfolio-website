package domain

import "time"

type Project struct {
	ID           string
	Name         string
	IDCategory   string
	IsDeployed   bool
	IsMaintained bool
	SourceCode   *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
