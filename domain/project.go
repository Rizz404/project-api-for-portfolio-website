package domain

import "time"

type Project struct {
	ID           string               `json:"id"`
	IDUser       string               `json:"id_user"`
	IDCategory   string               `json:"id_category"`
	IsDeployed   bool                 `json:"is_deployed"`
	IsMaintained bool                 `json:"is_maintained"`
	LiveDemo     *string              `json:"live_demo"`
	SourceCode   *string              `json:"source_code"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	User         User                 `json:"user"`
	Category     Category             `json:"category"`
	Translations []ProjectTranslation `json:"translations"`
	Techs        []Tech               `json:"techs"`
	Images       []ProjectImage       `json:"images"`
}

type CreateProjectPayload struct {
	IDCategory   string  `json:"id_category" form:"id_category" validate:"required,ulid"`
	IsDeployed   *bool   `json:"is_deployed" form:"is_deployed" validate:"required"`
	IsMaintained *bool   `json:"is_maintained" form:"is_maintained" validate:"required"`
	LiveDemo     *string `json:"live_demo,omitempty" form:"live_demo" validate:"omitempty,url"`
	SourceCode   *string `json:"source_code,omitempty" form:"source_code" validate:"omitempty,url"`
	// Translations []CreateProjectTranslationPayload `json:"translations" form:"translations" validate:"required,min=1,dive"`
	Techs []string `json:"techs" form:"techs" validate:"required,min=1,dive,ulid"`
}

type UpdateProjectPayload struct {
	IDCategory   *string  `json:"id_category,omitempty" form:"id_category" validate:"omitempty,ulid"`
	IsDeployed   *bool    `json:"is_deployed,omitempty" form:"is_deployed" validate:"omitempty"`
	IsMaintained *bool    `json:"is_maintained,omitempty" form:"is_maintained" validate:"omitempty"`
	LiveDemo     *string  `json:"live_demo,omitempty" form:"live_demo" validate:"omitempty,url"`
	SourceCode   *string  `json:"source_code,omitempty" form:"source_code" validate:"omitempty,url"`
	Techs        []string `json:"techs,omitempty" form:"techs" validate:"omitempty,min=1,dive,ulid"`
}
