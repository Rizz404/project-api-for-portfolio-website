package postgresql

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlTechStackRepository struct {
	queries *sqlc.Queries
}

func NewTechStackRepository(q *sqlc.Queries) *postgresqlTechStackRepository {
	return &postgresqlTechStackRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlTechStackRepository) CreateTechStack(ctx context.Context, payload *domain.TechStack) (domain.TechStack, error) {
	parsedIDProject, err := uuid.Parse(payload.IDProject)
	if err != nil {
		return domain.TechStack{}, err
	}
	parsedIDTech, err := uuid.Parse(payload.IDTech)
	if err != nil {
		return domain.TechStack{}, err
	}

	params := sqlc.CreateTechStackParams{
		IDProject: uuid.UUID(parsedIDProject),
		IDTech:    uuid.UUID(parsedIDTech),
	}

	sqlcTechStack, err := r.queries.CreateTechStack(ctx, params)
	if err != nil {
		return domain.TechStack{}, err
	}

	return toDomainTechStack(sqlcTechStack), nil
}

// *===========================UPDATE===========================*
func (r *postgresqlTechStackRepository) UpdateTechStack(ctx context.Context, payload *domain.TechStack) (domain.TechStack, error) {
	parsedIDProject, err := uuid.Parse(payload.IDProject)
	if err != nil {
		return domain.TechStack{}, err
	}
	parsedIDTech, err := uuid.Parse(payload.IDTech)
	if err != nil {
		return domain.TechStack{}, err
	}

	params := sqlc.UpdateTechStackParams{
		IDProject: uuid.UUID(parsedIDProject),
		IDTech:    uuid.UUID(parsedIDTech),
	}

	sqlcTechStack, err := r.queries.UpdateTechStack(ctx, params)
	if err != nil {
		return domain.TechStack{}, err
	}

	return toDomainTechStack(sqlcTechStack), nil
}

// *===========================DELETE===========================*
func (r *postgresqlTechStackRepository) DeleteTechStack(ctx context.Context, payload *domain.TechStack) error {
	parsedIDProject, err := uuid.Parse(payload.IDProject)
	if err != nil {
		return err
	}
	parsedIDTech, err := uuid.Parse(payload.IDTech)
	if err != nil {
		return err
	}

	params := sqlc.DeleteTechStackParams{
		IDProject: uuid.UUID(parsedIDProject),
		IDTech:    uuid.UUID(parsedIDTech),
	}

	err = r.queries.DeleteTechStack(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func toDomainTechStack(q sqlc.TechStack) domain.TechStack {
	domainIDProject := ulid.ULID(q.IDProject).String()
	domainIDTech := ulid.ULID(q.IDTech).String()

	return domain.TechStack{
		IDProject: domainIDProject,
		IDTech:    domainIDTech,
		CreatedAt: q.CreatedAt.Time,
		UpdatedAt: q.UpdatedAt.Time,
	}
}
