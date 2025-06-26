package language

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error)

	// * READ (MANY)
	GetLanguagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Language, error)
	GetLanguagesCursorFirst(ctx context.Context, limit int32) ([]domain.Language, error)
	SearchLanguages(ctx context.Context, searchTerm string) ([]domain.Language, error)
	SearchLanguagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Language, error)
	SearchLanguagesByName(ctx context.Context, name string) ([]domain.Language, error)
	SearchLanguagesByLangCode(ctx context.Context, langCode string) ([]domain.Language, error)

	// * READ (ONE & UTILITY)
	GetLanguage(ctx context.Context, id string) (domain.Language, error)
	GetLanguageByName(ctx context.Context, name string) (domain.Language, error)
	GetLanguageByLangCode(ctx context.Context, langCode string) (domain.Language, error)
	CheckLanguageExists(ctx context.Context, id string) (bool, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckLangCodeExists(ctx context.Context, langCode string) (bool, error)
	CountLanguages(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error)

	// * DELETE
	DeleteLanguage(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateLanguage(ctx context.Context, payload *domain.CreateLanguagePayload) (domain.Language, error) {

	params := domain.Language{
		ID:       ulid.Make().String(),
		Name:     payload.Name,
		LangCode: payload.LangCode,
	}

	nameExist, err := s.repo.CheckNameExists(ctx, params.Name)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}
	if nameExist {
		return domain.Language{}, domain.ErrConflict("Name already exists")
	}

	langCodeExist, err := s.repo.CheckLangCodeExists(ctx, params.LangCode)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}
	if langCodeExist {
		return domain.Language{}, domain.ErrConflict("LangCode already exists")
	}

	createdLanguage, err := s.repo.CreateLanguage(ctx, &params)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}

	return createdLanguage, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetLanguagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Language, error) {
	languages, err := s.repo.GetLanguagesPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

func (s *Service) GetLanguagesCursorFirst(ctx context.Context, limit int32) ([]domain.Language, error) {
	languages, err := s.repo.GetLanguagesCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

func (s *Service) SearchLanguages(ctx context.Context, searchTerm string) ([]domain.Language, error) {
	languages, err := s.repo.SearchLanguages(ctx, searchTerm)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

func (s *Service) SearchLanguagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Language, error) {
	languages, err := s.repo.SearchLanguagesPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

func (s *Service) SearchLanguagesByName(ctx context.Context, name string) ([]domain.Language, error) {
	languages, err := s.repo.SearchLanguagesByName(ctx, name)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

func (s *Service) SearchLanguagesByLangCode(ctx context.Context, langCode string) ([]domain.Language, error) {
	languages, err := s.repo.SearchLanguagesByLangCode(ctx, langCode)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return languages, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetLanguage(ctx context.Context, id string) (domain.Language, error) {
	language, err := s.repo.GetLanguage(ctx, id)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}

	return language, nil
}

func (s *Service) GetLanguageByName(ctx context.Context, name string) (domain.Language, error) {
	language, err := s.repo.GetLanguageByName(ctx, name)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}

	return language, nil
}

func (s *Service) GetLanguageByLangCode(ctx context.Context, langCode string) (domain.Language, error) {
	language, err := s.repo.GetLanguageByLangCode(ctx, langCode)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}

	return language, nil
}

func (s *Service) CheckLanguageExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckLanguageExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CheckLangCodeExists(ctx context.Context, langCode string) (bool, error) {
	exists, err := s.repo.CheckLangCodeExists(ctx, langCode)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CheckNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.CheckNameExists(ctx, name)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountLanguages(ctx context.Context) (int64, error) {
	count, err := s.repo.CountLanguages(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error) {
	exists, err := s.repo.CheckLanguageExists(ctx, payload.ID)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.Language{}, domain.ErrNotFound("Language")
	}

	if payload.Name != "" {
		nameExists, err := s.repo.CheckNameExists(ctx, payload.Name)
		if err != nil {
			return domain.Language{}, domain.ErrInternal(err)
		}
		if nameExists {
			currentLanguage, err := s.repo.GetLanguage(ctx, payload.ID)
			if err != nil {
				return domain.Language{}, domain.ErrInternal(err)
			}
			if currentLanguage.Name != payload.Name {
				return domain.Language{}, domain.ErrConflict("Name already exists")
			}
		}
	}

	if payload.LangCode != "" {
		langCodeExists, err := s.repo.CheckLangCodeExists(ctx, payload.LangCode)
		if err != nil {
			return domain.Language{}, domain.ErrInternal(err)
		}
		if langCodeExists {
			currentLanguage, err := s.repo.GetLanguage(ctx, payload.ID)
			if err != nil {
				return domain.Language{}, domain.ErrInternal(err)
			}
			if currentLanguage.LangCode != payload.LangCode {
				return domain.Language{}, domain.ErrConflict("LangCode already exists")
			}
		}
	}

	updatedLanguage, err := s.repo.UpdateLanguage(ctx, payload)
	if err != nil {
		return domain.Language{}, domain.ErrInternal(err)
	}

	return updatedLanguage, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteLanguage(ctx context.Context, id string) error {
	exists, err := s.repo.CheckLanguageExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("Language")
	}

	err = s.repo.DeleteLanguage(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
