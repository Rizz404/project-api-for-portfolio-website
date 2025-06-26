package domain

import "time"

type ProjectTranslation struct {
	ID          string
	IDProject   string
	Description string
	LangCode    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
