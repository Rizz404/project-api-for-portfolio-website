package project

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateProject(ctx context.Context, payload *domain.Project) (domain.Project, error)

	// * READ (MANY)
	GetProjectsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Project, error)
	GetProjectsCursorFirst(ctx context.Context, limit int32) ([]domain.Project, error)
	SearchProjectsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Project, error)

	// * READ (ONE & UTILITY)
	GetProject(ctx context.Context, id string) (domain.Project, error)
	GetProjectByTranslatedName(ctx context.Context, name string) (domain.Project, error)
	CheckProjectExists(ctx context.Context, id string) (bool, error)
	CountProjects(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProject(ctx context.Context, payload *domain.Project) (domain.Project, error)

	// * DELETE
	DeleteProject(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateProject(ctx context.Context, payload *domain.CreateProjectPayload, userId string) (domain.Project, error) {

	params := domain.Project{
		ID:         ulid.Make().String(),
		IDUser:     userId,
		IDCategory: payload.IDCategory,
		// IsDeployed:   payload.IsDeployed,
		// IsMaintained: payload.IsMaintained,
		LiveDemo:   payload.LiveDemo,
		SourceCode: payload.SourceCode,
	}

	createdProject, err := s.repo.CreateProject(ctx, &params)
	if err != nil {
		return domain.Project{}, domain.ErrInternal(err)
	}

	return createdProject, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetProjectsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Project, error) {
	projects, err := s.repo.GetProjectsPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projects, nil
}

func (s *Service) GetProjectsCursorFirst(ctx context.Context, limit int32) ([]domain.Project, error) {
	projects, err := s.repo.GetProjectsCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projects, nil
}

func (s *Service) SearchProjectsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Project, error) {
	projects, err := s.repo.SearchProjectsPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projects, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetProject(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.repo.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, domain.ErrInternal(err)
	}

	return project, nil
}

func (s *Service) GetProjectByTranslatedName(ctx context.Context, name string) (domain.Project, error) {
	project, err := s.repo.GetProjectByTranslatedName(ctx, name)
	if err != nil {
		return domain.Project{}, domain.ErrInternal(err)
	}

	return project, nil
}

func (s *Service) CheckProjectExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckProjectExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountProjects(ctx context.Context) (int64, error) {
	count, err := s.repo.CountProjects(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateProject(ctx context.Context, payload *domain.Project) (domain.Project, error) {
	exists, err := s.repo.CheckProjectExists(ctx, payload.ID)
	if err != nil {
		return domain.Project{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.Project{}, domain.ErrNotFound("Project")
	}

	updatedProject, err := s.repo.UpdateProject(ctx, payload)
	if err != nil {
		return domain.Project{}, domain.ErrInternal(err)
	}

	return updatedProject, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteProject(ctx context.Context, id string) error {
	exists, err := s.repo.CheckProjectExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("Project")
	}

	err = s.repo.DeleteProject(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
