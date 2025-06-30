package projectimage

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error)
	CreateProjectImagesBatch(ctx context.Context, payload *[]domain.ProjectImage) ([]domain.ProjectImage, error)

	// * READ (MANY)
	GetProjectImagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectImage, error)
	GetProjectImagesCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectImage, error)
	SearchProjectImages(ctx context.Context, searchTerm string) ([]domain.ProjectImage, error)
	SearchProjectImagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectImage, error)

	// * READ (ONE & UTILITY)
	GetProjectImage(ctx context.Context, id string) (domain.ProjectImage, error)
	CheckProjectImageExists(ctx context.Context, id string) (bool, error)
	CountProjectImages(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error)

	// * DELETE
	DeleteProjectImage(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateProjectImage(ctx context.Context, payload *domain.CreateProjectImagePayload) (domain.ProjectImage, error) {

	params := domain.ProjectImage{
		ID:        ulid.Make().String(),
		IDProject: payload.IDProject,
		FileName:  payload.FileName,
		Url:       payload.Url,
	}

	createdProjectImage, err := s.repo.CreateProjectImage(ctx, &params)
	if err != nil {
		return domain.ProjectImage{}, domain.ErrInternal(err)
	}

	return createdProjectImage, nil
}

func (s *Service) CreateProjectImagesBatch(ctx context.Context, payload *[]domain.CreateProjectImagePayload) ([]domain.ProjectImage, error) {
	batchItems := make([]domain.ProjectImage, len(*payload))

	for _, projectImage := range batchItems {
		params := domain.ProjectImage{
			ID:        ulid.Make().String(),
			IDProject: projectImage.IDProject,
			FileName:  projectImage.FileName,
			Url:       projectImage.Url,
		}

		batchItems = append(batchItems, params)
	}

	createdProjectImage, err := s.repo.CreateProjectImagesBatch(ctx, &batchItems)
	if err != nil {
		return []domain.ProjectImage{}, domain.ErrInternal(err)
	}

	return createdProjectImage, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetProjectImagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectImage, error) {
	projectImages, err := s.repo.GetProjectImagesPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectImages, nil
}

func (s *Service) GetProjectImagesCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectImage, error) {
	projectImages, err := s.repo.GetProjectImagesCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectImages, nil
}

func (s *Service) SearchProjectImages(ctx context.Context, searchTerm string) ([]domain.ProjectImage, error) {
	projectImages, err := s.repo.SearchProjectImages(ctx, searchTerm)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectImages, nil
}

func (s *Service) SearchProjectImagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectImage, error) {
	projectImages, err := s.repo.SearchProjectImagesPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return projectImages, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetProjectImage(ctx context.Context, id string) (domain.ProjectImage, error) {
	projectImage, err := s.repo.GetProjectImage(ctx, id)
	if err != nil {
		return domain.ProjectImage{}, domain.ErrInternal(err)
	}

	return projectImage, nil
}

func (s *Service) CheckProjectImageExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckProjectImageExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountProjectImages(ctx context.Context) (int64, error) {
	count, err := s.repo.CountProjectImages(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error) {
	exists, err := s.repo.CheckProjectImageExists(ctx, payload.ID)
	if err != nil {
		return domain.ProjectImage{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.ProjectImage{}, domain.ErrNotFound("ProjectImage")
	}

	updatedProjectImage, err := s.repo.UpdateProjectImage(ctx, payload)
	if err != nil {
		return domain.ProjectImage{}, domain.ErrInternal(err)
	}

	return updatedProjectImage, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteProjectImage(ctx context.Context, id string) error {
	exists, err := s.repo.CheckProjectImageExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("ProjectImage")
	}

	err = s.repo.DeleteProjectImage(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
