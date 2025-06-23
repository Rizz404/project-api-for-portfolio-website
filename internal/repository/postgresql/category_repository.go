package postgresql

import (
	"context"
	"database/sql"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlCategoryRepository struct {
	queries *sqlc.Queries
}

func NewCategoryRepository(q *sqlc.Queries) *postgresqlCategoryRepository {
	return &postgresqlCategoryRepository{
		queries: q,
	}
}

func (r *postgresqlCategoryRepository) CreateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Category{}, err
	}

	var sqlcDescription sql.NullString
	if payload.Description != nil {
		sqlcDescription.String = *payload.Description
		sqlcDescription.Valid = true
	}

	sqlcCategory, err := r.queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		ID:          uuid.UUID(parsedID),
		Name:        payload.Name,
		Description: sqlcDescription,
	})
	if err != nil {
		return domain.Category{}, err
	}

	return toDomainCategory(sqlcCategory), nil
}

func (r *postgresqlCategoryRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	sqlcCategories, err := r.queries.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainCategories(sqlcCategories), nil
}

func (r *postgresqlCategoryRepository) GetCategory(ctx context.Context, id string) (domain.Category, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.Category{}, err
	}

	sqlcCategory, err := r.queries.GetCategory(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.Category{}, err
	}

	return toDomainCategory(sqlcCategory), nil
}

func (r *postgresqlCategoryRepository) UpdateCategory(ctx context.Context, payload *domain.Category) (domain.Category, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Category{}, err
	}

	var sqlcDescription sql.NullString
	if payload.Description != nil {
		sqlcDescription.String = *payload.Description
		sqlcDescription.Valid = true
	}

	sqlcCategory, err := r.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:          uuid.UUID(parsedID),
		Name:        payload.Name,
		Description: sqlcDescription,
	})
	if err != nil {
		return domain.Category{}, err
	}

	return toDomainCategory(sqlcCategory), nil
}

func (r *postgresqlCategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteCategory(ctx, uuid.UUID(parsedID))
}

func toDomainCategory(q sqlc.Category) domain.Category {
	var domainDesc *string
	if q.Description.Valid {
		domainDesc = &q.Description.String
	}

	domainID := ulid.ULID(q.ID).String()

	return domain.Category{
		ID:          domainID,
		Name:        q.Name,
		Description: domainDesc,
		CreatedAt:   q.CreatedAt.Time,
		UpdatedAt:   q.UpdatedAt.Time,
	}
}

func toDomainCategories(q []sqlc.Category) []domain.Category {
	domainCategories := make([]domain.Category, len(q))

	for i, sqlCat := range q {
		domainCategories[i] = toDomainCategory(sqlCat)
	}
	return domainCategories
}
