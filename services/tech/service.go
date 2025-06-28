package tech

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error)

	// * READ (MANY)
	GetTechsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Tech, error)
	GetTechsCursorFirst(ctx context.Context, limit int32) ([]domain.Tech, error)
	SearchTechs(ctx context.Context, searchTerm string) ([]domain.Tech, error)
	SearchTechsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Tech, error)

	// * READ (ONE & UTILITY)
	GetTech(ctx context.Context, id string) (domain.Tech, error)
	GetTechByName(ctx context.Context, name string) (domain.Tech, error)
	CheckTechExists(ctx context.Context, id string) (bool, error)
	CheckTechNameExists(ctx context.Context, name string) (bool, error)
	CountTechs(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error)

	// * DELETE
	DeleteTech(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateTech(ctx context.Context, payload *domain.CreateTechPayload) (domain.Tech, error) {

	params := domain.Tech{
		ID:          ulid.Make().String(),
		Name:        payload.Name,
		Description: payload.Description,
		LogoURL:     payload.LogoURL,
	}

	nameExist, err := s.repo.CheckTechNameExists(ctx, params.Name)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}
	if nameExist {
		return domain.Tech{}, domain.ErrConflict("Name already exists")
	}

	createdTech, err := s.repo.CreateTech(ctx, &params)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}

	return createdTech, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetTechsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Tech, error) {
	techs, err := s.repo.GetTechsPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return techs, nil
}

func (s *Service) GetTechsCursorFirst(ctx context.Context, limit int32) ([]domain.Tech, error) {
	techs, err := s.repo.GetTechsCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return techs, nil
}

func (s *Service) SearchTechs(ctx context.Context, searchTerm string) ([]domain.Tech, error) {
	techs, err := s.repo.SearchTechs(ctx, searchTerm)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return techs, nil
}

func (s *Service) SearchTechsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Tech, error) {
	techs, err := s.repo.SearchTechsPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return techs, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetTech(ctx context.Context, id string) (domain.Tech, error) {
	tech, err := s.repo.GetTech(ctx, id)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}

	return tech, nil
}

func (s *Service) GetTechByName(ctx context.Context, name string) (domain.Tech, error) {
	tech, err := s.repo.GetTechByName(ctx, name)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}

	return tech, nil
}

func (s *Service) CheckTechExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckTechExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CheckTechNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.CheckTechNameExists(ctx, name)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountTechs(ctx context.Context) (int64, error) {
	count, err := s.repo.CountTechs(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error) {
	exists, err := s.repo.CheckTechExists(ctx, payload.ID)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.Tech{}, domain.ErrNotFound("Tech")
	}

	if payload.Name != "" {
		nameExists, err := s.repo.CheckTechNameExists(ctx, payload.Name)
		if err != nil {
			return domain.Tech{}, domain.ErrInternal(err)
		}
		if nameExists {
			currentTech, err := s.repo.GetTech(ctx, payload.ID)
			if err != nil {
				return domain.Tech{}, domain.ErrInternal(err)
			}
			if currentTech.Name != payload.Name {
				return domain.Tech{}, domain.ErrConflict("Name already exists")
			}
		}
	}

	updatedTech, err := s.repo.UpdateTech(ctx, payload)
	if err != nil {
		return domain.Tech{}, domain.ErrInternal(err)
	}

	return updatedTech, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteTech(ctx context.Context, id string) error {
	exists, err := s.repo.CheckTechExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("Tech")
	}

	err = s.repo.DeleteTech(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
