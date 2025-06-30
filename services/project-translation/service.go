package projecttranslation

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error)
	CreateProjectTranslationsBatch(ctx context.Context, payload *[]domain.ProjectTranslation) ([]domain.ProjectTranslation, error)

	// * READ (MANY)
	GetProjectTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectTranslation, error)
	GetProjectTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectTranslation, error)
	SearchProjectTranslations(ctx context.Context, searchTerm string) ([]domain.ProjectTranslation, error)
	SearchProjectTranslationsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectTranslation, error)

	// * READ (ONE & UTILITY)
	GetProjectTranslation(ctx context.Context, id string) (domain.ProjectTranslation, error)
	GetProjectTranslationByName(ctx context.Context, name string) (domain.ProjectTranslation, error)
	CheckProjectTranslationExists(ctx context.Context, id string) (bool, error)
	CountProjectTranslations(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error)

	// * DELETE
	DeleteProjectTranslation(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateProjectTranslation(ctx context.Context, payload *domain.CreateProjectTranslationPayload) (domain.ProjectTranslation, error) {

	params := domain.ProjectTranslation{
		ID:          ulid.Make().String(),
		IDProject:   payload.IDProject,
		IDLanguage:  payload.IDLanguage,
		Name:        payload.Name,
		Description: payload.Description,
	}

	createdProjectTranslation, err := s.repo.CreateProjectTranslation(ctx, &params)
	if err != nil {
		return domain.ProjectTranslation{}, domain.ErrInternal(err)
	}

	return createdProjectTranslation, nil
}

func (s *Service) CreateProjectTranslationsBatch(ctx context.Context, payload *[]domain.CreateProjectTranslationPayload) ([]domain.ProjectTranslation, error) {
	batchItems := make([]domain.ProjectTranslation, len(*payload))

	for _, projectTranslation := range batchItems {
		params := domain.ProjectTranslation{
			ID:          ulid.Make().String(),
			IDProject:   projectTranslation.IDProject,
			Name:        projectTranslation.Name,
			Description: projectTranslation.Description,
		}

		batchItems = append(batchItems, params)
	}

	createdProjectTranslation, err := s.repo.CreateProjectTranslationsBatch(ctx, &batchItems)
	if err != nil {
		return []domain.ProjectTranslation{}, domain.ErrInternal(err)
	}

	return createdProjectTranslation, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetProjectTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectTranslation, error) {
	projectTranslations, err := s.repo.GetProjectTranslationsPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectTranslations, nil
}

func (s *Service) GetProjectTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectTranslation, error) {
	projectTranslations, err := s.repo.GetProjectTranslationsCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectTranslations, nil
}

func (s *Service) SearchProjectTranslations(ctx context.Context, searchTerm string) ([]domain.ProjectTranslation, error) {
	projectTranslations, err := s.repo.SearchProjectTranslations(ctx, searchTerm)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectTranslations, nil
}

func (s *Service) SearchProjectTranslationsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectTranslation, error) {
	projectTranslations, err := s.repo.SearchProjectTranslationsPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectTranslations, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetProjectTranslation(ctx context.Context, id string) (domain.ProjectTranslation, error) {
	projectTranslation, err := s.repo.GetProjectTranslation(ctx, id)
	if err != nil {
		return domain.ProjectTranslation{}, domain.ErrInternal(err)
	}

	return projectTranslation, nil
}

func (s *Service) GetProjectTranslationByName(ctx context.Context, name string) (domain.ProjectTranslation, error) {
	projectTranslation, err := s.repo.GetProjectTranslationByName(ctx, name)
	if err != nil {
		return domain.ProjectTranslation{}, domain.ErrInternal(err)
	}

	return projectTranslation, nil
}

func (s *Service) CheckProjectTranslationExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckProjectTranslationExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountProjectTranslations(ctx context.Context) (int64, error) {
	count, err := s.repo.CountProjectTranslations(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error) {
	exists, err := s.repo.CheckProjectTranslationExists(ctx, payload.ID)
	if err != nil {
		return domain.ProjectTranslation{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.ProjectTranslation{}, domain.ErrNotFound("ProjectTranslation")
	}

	updatedProjectTranslation, err := s.repo.UpdateProjectTranslation(ctx, payload)
	if err != nil {
		return domain.ProjectTranslation{}, domain.ErrInternal(err)
	}

	return updatedProjectTranslation, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteProjectTranslation(ctx context.Context, id string) error {
	exists, err := s.repo.CheckProjectTranslationExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("ProjectTranslation")
	}

	err = s.repo.DeleteProjectTranslation(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
