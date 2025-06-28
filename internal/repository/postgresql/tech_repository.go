package postgresql

import (
	"context"
	"database/sql"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlTechRepository struct {
	queries *sqlc.Queries
}

func NewTechRepository(q *sqlc.Queries) *postgresqlTechRepository {
	return &postgresqlTechRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlTechRepository) CreateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Tech{}, err
	}

	params := sqlc.CreateTechParams{
		ID:   uuid.UUID(parsedID),
		Name: payload.Name,
	}

	if payload.Description != nil {
		params.Description = sql.NullString{String: *payload.Description, Valid: true}
	}

	if payload.LogoURL != nil {
		params.LogoUrl = sql.NullString{String: *payload.LogoURL, Valid: true}
	}

	sqlcTech, err := r.queries.CreateTech(ctx, params)
	if err != nil {
		return domain.Tech{}, err
	}

	return toDomainTech(sqlcTech), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlTechRepository) GetTechsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Tech, error) {
	sqlcTechs, err := r.queries.GetTechsPaginated(ctx, sqlc.GetTechsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainTechs(sqlcTechs), nil
}

func (r *postgresqlTechRepository) GetTechsCursorFirst(ctx context.Context, limit int32) ([]domain.Tech, error) {
	sqlcTechs, err := r.queries.GetTechsCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainTechs(sqlcTechs), nil
}

func (r *postgresqlTechRepository) SearchTechs(ctx context.Context, searchTerm string) ([]domain.Tech, error) {
	sqlcTechs, err := r.queries.SearchTechs(ctx, sql.NullString{String: searchTerm, Valid: searchTerm != ""})
	if err != nil {
		return nil, err
	}
	return toDomainTechs(sqlcTechs), nil
}

func (r *postgresqlTechRepository) SearchTechsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Tech, error) {
	sqlcTechs, err := r.queries.SearchTechsPaginated(ctx, sqlc.SearchTechsPaginatedParams{
		Column1: sql.NullString{String: searchTerm, Valid: searchTerm != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainTechs(sqlcTechs), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlTechRepository) GetTech(ctx context.Context, id string) (domain.Tech, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.Tech{}, err
	}
	sqlcTech, err := r.queries.GetTech(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.Tech{}, err
	}
	return toDomainTech(sqlcTech), nil
}

func (r *postgresqlTechRepository) GetTechByName(ctx context.Context, name string) (domain.Tech, error) {
	sqlcTech, err := r.queries.GetTechByName(ctx, name)
	if err != nil {
		return domain.Tech{}, err
	}
	return toDomainTech(sqlcTech), nil
}

func (r *postgresqlTechRepository) CheckTechExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckTechExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlTechRepository) CheckTechNameExists(ctx context.Context, name string) (bool, error) {
	return r.queries.CheckTechNameExists(ctx, name)
}

func (r *postgresqlTechRepository) CountTechs(ctx context.Context) (int64, error) {
	return r.queries.CountTechs(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlTechRepository) UpdateTech(ctx context.Context, payload *domain.Tech) (domain.Tech, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Tech{}, err
	}

	params := sqlc.UpdateTechParams{
		ID:   uuid.UUID(parsedID),
		Name: payload.Name,
	}

	if payload.Description != nil {
		params.Description = sql.NullString{String: *payload.Description, Valid: true}
	}

	if payload.LogoURL != nil {
		params.LogoUrl = sql.NullString{String: *payload.LogoURL, Valid: true}
	}

	sqlcTech, err := r.queries.UpdateTech(ctx, params)
	if err != nil {
		return domain.Tech{}, err
	}

	return toDomainTech(sqlcTech), nil
}

// *===========================DELETE===========================*
func (r *postgresqlTechRepository) DeleteTech(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteTech(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainTech(q sqlc.Tech) domain.Tech {
	domainID := ulid.ULID(q.ID).String()

	var domainDescription *string
	if q.Description.Valid {
		domainDescription = &q.Description.String
	}

	return domain.Tech{
		ID:          domainID,
		Name:        q.Name,
		Description: domainDescription,
		CreatedAt:   q.CreatedAt.Time,
		UpdatedAt:   q.UpdatedAt.Time,
	}
}

func toDomainTechs(q []sqlc.Tech) []domain.Tech {
	domainTechs := make([]domain.Tech, len(q))

	for i, sqlcTech := range q {
		domainTechs[i] = toDomainTech(sqlcTech)
	}
	return domainTechs
}
