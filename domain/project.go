package domain

import "time"

type Project struct {
	Id           string
	Name         string
	Description  *string
	IdCategory   string
	IsDeployed   bool
	IsMaintained bool
	SourceCode   *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
