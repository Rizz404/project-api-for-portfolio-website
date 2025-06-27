package usertranslation

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateUserTranslation(ctx context.Context, payload *domain.UserTranslation) (domain.UserTranslation, error)

	// * READ (MANY)
	GetUserTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.UserTranslation, error)
	GetUserTranslationsByUserIDPaginated(ctx context.Context, idUser string, limit int32, offset int32) ([]domain.UserTranslation, error)
	GetUserTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.UserTranslation, error)

	// * READ (ONE & UTILITY)
	GetUserTranslation(ctx context.Context, id string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangID(ctx context.Context, idUser string, idLanguage string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangName(ctx context.Context, idUser string, langName string) (domain.UserTranslation, error)
	GetUserTranslationByUserIDAndLangCode(ctx context.Context, idUser string, langCode string) (domain.UserTranslation, error)
	CheckUserTranslationExists(ctx context.Context, id string) (bool, error)
	CountUserTranslations(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateUserTranslation(ctx context.Context, payload *domain.UserTranslation) (domain.UserTranslation, error)

	// * DELETE
	DeleteUserTranslation(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
	// Kamu mungkin perlu menambahkan repository lain di sini jika ada dependensi
	// contoh: userRepo user.Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateUserTranslation(ctx context.Context, payload *domain.CreateUserTranslationPayload) (domain.UserTranslation, error) {
	_, err := s.repo.GetUserTranslationByUserIDAndLangID(ctx, payload.IDUser, payload.IDLanguage)
	if err == nil {
		return domain.UserTranslation{}, domain.ErrConflict("User translation for this language already exists")
	}

	if err != domain.ErrNotFound("UserTranslation") { // Asumsi ErrNotFound memiliki format ini
		return domain.UserTranslation{}, domain.ErrInternal(err)
	}

	params := domain.UserTranslation{
		ID:               ulid.Make().String(),
		IDUser:           payload.IDUser,
		IDLanguage:       payload.IDLanguage,
		Bio:              payload.Bio,
		AboutMe:          payload.AboutMe,
		AdditionalSkills: payload.AdditionalSkills,
		Languages:        payload.Languages,
		Quote:            payload.Quote,
	}

	createdUserTranslation, err := s.repo.CreateUserTranslation(ctx, &params)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrInternal(err)
	}

	return createdUserTranslation, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetUserTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.UserTranslation, error) {
	translations, err := s.repo.GetUserTranslationsPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}
	return translations, nil
}

func (s *Service) GetUserTranslationsByUserIDPaginated(ctx context.Context, idUser string, limit int32, offset int32) ([]domain.UserTranslation, error) {
	translations, err := s.repo.GetUserTranslationsByUserIDPaginated(ctx, idUser, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}
	return translations, nil
}

func (s *Service) GetUserTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.UserTranslation, error) {
	translations, err := s.repo.GetUserTranslationsCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}
	return translations, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetUserTranslation(ctx context.Context, id string) (domain.UserTranslation, error) {
	translation, err := s.repo.GetUserTranslation(ctx, id)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrNotFound("UserTranslation")
	}
	return translation, nil
}

func (s *Service) GetUserTranslationByUserIDAndLangID(ctx context.Context, idUser string, idLanguage string) (domain.UserTranslation, error) {
	translation, err := s.repo.GetUserTranslationByUserIDAndLangID(ctx, idUser, idLanguage)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrNotFound("UserTranslation")
	}
	return translation, nil
}

func (s *Service) GetUserTranslationByUserIDAndLangName(ctx context.Context, idUser string, langName string) (domain.UserTranslation, error) {
	translation, err := s.repo.GetUserTranslationByUserIDAndLangName(ctx, idUser, langName)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrNotFound("UserTranslation")
	}
	return translation, nil
}

func (s *Service) GetUserTranslationByUserIDAndLangCode(ctx context.Context, idUser string, langCode string) (domain.UserTranslation, error) {
	translation, err := s.repo.GetUserTranslationByUserIDAndLangCode(ctx, idUser, langCode)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrNotFound("UserTranslation")
	}
	return translation, nil
}

func (s *Service) CheckUserTranslationExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckUserTranslationExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountUserTranslations(ctx context.Context) (int64, error) {
	count, err := s.repo.CountUserTranslations(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateUserTranslation(ctx context.Context, id string, payload *domain.UpdateUserTranslationPayload) (domain.UserTranslation, error) {
	exists, err := s.repo.CheckUserTranslationExists(ctx, id)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.UserTranslation{}, domain.ErrNotFound("UserTranslation")
	}

	currentTranslation, err := s.repo.GetUserTranslation(ctx, id)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrInternal(err)
	}

	currentTranslation.Bio = payload.Bio
	currentTranslation.AboutMe = payload.AboutMe
	currentTranslation.AdditionalSkills = nil
	if payload.AdditionalSkills != nil {
		currentTranslation.AdditionalSkills = *payload.AdditionalSkills
	}
	if payload.Languages != nil {
		currentTranslation.Languages = *payload.Languages
	}
	currentTranslation.Quote = payload.Quote

	updatedTranslation, err := s.repo.UpdateUserTranslation(ctx, &currentTranslation)
	if err != nil {
		return domain.UserTranslation{}, domain.ErrInternal(err)
	}

	return updatedTranslation, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteUserTranslation(ctx context.Context, id string) error {
	exists, err := s.repo.CheckUserTranslationExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("UserTranslation")
	}

	err = s.repo.DeleteUserTranslation(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
