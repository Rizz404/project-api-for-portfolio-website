package category

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/oklog/ulid/v2"
)

// * Repository adalah kontrak/interface yang harus dipenuhi oleh lapisan database.
// * Signature-nya harus cocok dengan apa yang bisa dilakukan oleh repository.
type Repository interface {
	CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	GetCategory(ctx context.Context, id string) (domain.Category, error)
	UpdateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

// * Service adalah struct yang berisi logika bisnis.
type Service struct {
	repo Repository
}

// * NewService adalah constructor untuk Category Service.
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateCategory(ctx context.Context, payload *domain.CreateCategoryPayload) (domain.Category, error) {
	newID := ulid.Make().String()

	categoryToStore := &domain.Category{
		ID:          newID,
		Name:        payload.Name,
		Description: payload.Description,
	}

	return s.repo.CreateCategory(ctx, categoryToStore)
}

func (s *Service) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) GetCategory(ctx context.Context, id string) (domain.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *Service) UpdateCategory(ctx context.Context, id string, payload *domain.UpdateCategoryPayload) (domain.Category, error) {
	existingCategory, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		return domain.Category{}, err
	}

	if payload.Name != nil {
		existingCategory.Name = *payload.Name
	}
	if payload.Description != nil {
		existingCategory.Description = payload.Description
	}

	return s.repo.UpdateCategory(ctx, &existingCategory)
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteCategory(ctx, id)
}
