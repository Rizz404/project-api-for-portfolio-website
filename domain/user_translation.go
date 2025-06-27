package domain

import "time"

type UserTranslation struct {
	ID               string    `json:"id"`
	IDUser           string    `json:"id_user"`
	IDLanguage       string    `json:"id_language"`
	Bio              *string   `json:"bio,omitempty"`
	AboutMe          *string   `json:"about_me,omitempty"`
	AdditionalSkills []string  `json:"additional_skills"`
	Languages        []string  `json:"languages"`
	Quote            *string   `json:"quote,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateUserTranslationPayload struct {
	IDUser           string   `json:"id_user" form:"id_user" validate:"required,ulid"`
	IDLanguage       string   `json:"id_language" form:"id_language" validate:"required,ulid"`
	Bio              *string  `json:"bio,omitempty" form:"bio" validate:"omitempty,min=10,max=500"`
	AboutMe          *string  `json:"about_me,omitempty" form:"about_me" validate:"omitempty,min=20,max=2000"`
	AdditionalSkills []string `json:"additional_skills" form:"additional_skills" validate:"required,min=1,dive,min=2,max=50"`
	Languages        []string `json:"languages" form:"languages" validate:"required,min=1,dive,min=2,max=50"`
	Quote            *string  `json:"quote,omitempty" form:"quote" validate:"omitempty,max=255"`
}

type UpdateUserTranslationPayload struct {
	Bio              *string   `json:"bio,omitempty" validate:"omitempty,min=10,max=500"`
	AboutMe          *string   `json:"about_me,omitempty" validate:"omitempty,min=20,max=2000"`
	AdditionalSkills *[]string `json:"additional_skills,omitempty" validate:"omitempty,min=1,dive,min=2,max=50"`
	Languages        *[]string `json:"languages,omitempty" validate:"omitempty,min=1,dive,min=2,max=50"`
	Quote            *string   `json:"quote,omitempty" validate:"omitempty,max=255"`
}
