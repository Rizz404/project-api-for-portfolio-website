package techstack

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
)

type Repository interface {
	// * CREATE
	CreateTechStack(ctx context.Context, payload *domain.TechStack) (domain.TechStack, error)

	// * UPDATE
	UpdateTechStack(ctx context.Context, payload *domain.TechStack) (domain.TechStack, error)

	// * DELETE
	DeleteTechStack(ctx context.Context, payload *domain.TechStack) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateTechStack(ctx context.Context, payload *domain.CreateTechStackPayload) (domain.TechStack, error) {
	params := domain.TechStack{
		IDProject: payload.IDProject,
		IDTech:    payload.IDTech,
	}

	createdTechStack, err := s.repo.CreateTechStack(ctx, &params)
	if err != nil {
		return domain.TechStack{}, domain.ErrInternal(err)
	}

	return createdTechStack, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateTechStack(ctx context.Context, payload *domain.UpdateTechStackPayload) (domain.TechStack, error) {
	params := domain.TechStack{}

	if payload.IDProject != nil {
		params.IDProject = *payload.IDProject
	}

	if payload.IDTech != nil {
		params.IDTech = *payload.IDTech
	}

	createdTechStack, err := s.repo.UpdateTechStack(ctx, &params)
	if err != nil {
		return domain.TechStack{}, domain.ErrInternal(err)
	}

	return createdTechStack, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteTechStack(ctx context.Context, payload *domain.DeleteTechStackPayload) error {
	params := domain.TechStack{
		IDProject: payload.IDProject,
		IDTech:    payload.IDTech,
	}

	err := s.repo.DeleteTechStack(ctx, &params)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
