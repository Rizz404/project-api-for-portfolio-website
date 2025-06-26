package domain

import "time"

type UserTranslation struct {
	ID               string    `json:"id"`
	IDUser           string    `json:"id_user"`
	Bio              *string   `json:"bio,omitempty"`
	AboutMe          *string   `json:"about_me,omitempty"`
	AdditionalSkills []string  `json:"additional_skills"`
	Languages        []string  `json:"languages"`
	Quote            *string   `json:"quote,omitempty"`
	LangCode         string    `json:"lang_code"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
