package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlProjectTranslationRepository struct {
	queries *sqlc.Queries
}

func NewProjectTranslationRepository(q *sqlc.Queries) *postgresqlProjectTranslationRepository {
	return &postgresqlProjectTranslationRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlProjectTranslationRepository) CreateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	parsedIDProject, err := ulid.Parse(payload.IDProject)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	parsedIDLanguage, err := ulid.Parse(payload.IDLanguage)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}

	params := sqlc.CreateProjectTranslationParams{
		ID:         uuid.UUID(parsedID),
		IDProject:  uuid.UUID(parsedIDProject),
		IDLanguage: uuid.UUID(parsedIDLanguage),
		Name:       payload.Name,
	}

	if payload.Description != nil {
		params.Description = sql.NullString{String: *payload.Description, Valid: true}
	}

	sqlcProjectTranslation, err := r.queries.CreateProjectTranslation(ctx, params)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}

	return toDomainProjectTranslation(sqlcProjectTranslation), nil
}

func (r *postgresqlProjectTranslationRepository) CreateProjectTranslationsBatch(ctx context.Context, payload *[]domain.ProjectTranslation) ([]domain.ProjectTranslation, error) {

	batchItems := make([]domain.ProjectTranslation, len(*payload))

	for _, pi := range batchItems {
		parsedID, err := ulid.Parse(pi.ID)
		if err != nil {
			return nil, err
		}
		parsedIDProject, err := ulid.Parse(pi.IDProject)
		if err != nil {
			return nil, err
		}
		parsedIDLanguage, err := ulid.Parse(pi.IDLanguage)
		if err != nil {
			return nil, err
		}

		batchItems = append(batchItems, domain.ProjectTranslation{
			ID:          uuid.UUID(parsedID).String(),
			IDProject:   uuid.UUID(parsedIDProject).String(),
			IDLanguage:  uuid.UUID(parsedIDLanguage).String(),
			Name:        pi.Name,
			Description: pi.Description,
		})
	}

	jsonData, err := json.Marshal(batchItems)
	if err != nil {
		return nil, err
	}

	sqlcProjectTranslations, err := r.queries.CreateProjectTranslationsBatch(ctx, json.RawMessage(jsonData))
	if err != nil {
		return nil, err
	}

	return toDomainProjectTranslations(sqlcProjectTranslations), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlProjectTranslationRepository) GetProjectTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectTranslation, error) {
	sqlcProjectTranslations, err := r.queries.GetProjectTranslationsPaginated(ctx, sqlc.GetProjectTranslationsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectTranslations(sqlcProjectTranslations), nil
}

func (r *postgresqlProjectTranslationRepository) GetProjectTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectTranslation, error) {
	sqlcProjectTranslations, err := r.queries.GetProjectTranslationsCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainProjectTranslations(sqlcProjectTranslations), nil
}

func (r *postgresqlProjectTranslationRepository) SearchProjectTranslations(ctx context.Context, searchTerm string) ([]domain.ProjectTranslation, error) {
	sqlcProjectTranslations, err := r.queries.SearchProjectTranslations(ctx, sql.NullString{String: searchTerm, Valid: searchTerm != ""})
	if err != nil {
		return nil, err
	}
	return toDomainProjectTranslations(sqlcProjectTranslations), nil
}

func (r *postgresqlProjectTranslationRepository) SearchProjectTranslationsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectTranslation, error) {
	sqlcProjectTranslations, err := r.queries.SearchProjectTranslationsPaginated(ctx, sqlc.SearchProjectTranslationsPaginatedParams{
		Column1: sql.NullString{String: searchTerm, Valid: searchTerm != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectTranslations(sqlcProjectTranslations), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlProjectTranslationRepository) GetProjectTranslation(ctx context.Context, id string) (domain.ProjectTranslation, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	sqlcProjectTranslation, err := r.queries.GetProjectTranslation(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	return toDomainProjectTranslation(sqlcProjectTranslation), nil
}

func (r *postgresqlProjectTranslationRepository) GetProjectTranslationByName(ctx context.Context, name string) (domain.ProjectTranslation, error) {
	sqlcProjectTranslation, err := r.queries.GetProjectTranslationByName(ctx, name)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	return toDomainProjectTranslation(sqlcProjectTranslation), nil
}

func (r *postgresqlProjectTranslationRepository) CheckProjectTranslationExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckProjectTranslationExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlProjectTranslationRepository) CountProjectTranslations(ctx context.Context) (int64, error) {
	return r.queries.CountProjectTranslations(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlProjectTranslationRepository) UpdateProjectTranslation(ctx context.Context, payload *domain.ProjectTranslation) (domain.ProjectTranslation, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}
	parsedIDLanguage, err := ulid.Parse(payload.IDLanguage)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}

	params := sqlc.UpdateProjectTranslationParams{
		ID:         uuid.UUID(parsedID),
		IDLanguage: uuid.UUID(parsedIDLanguage),
		Name:       payload.Name,
	}

	if payload.Description != nil {
		params.Description = sql.NullString{String: *payload.Description, Valid: true}
	}

	sqlcProjectTranslation, err := r.queries.UpdateProjectTranslation(ctx, params)
	if err != nil {
		return domain.ProjectTranslation{}, err
	}

	return toDomainProjectTranslation(sqlcProjectTranslation), nil
}

// *===========================DELETE===========================*
func (r *postgresqlProjectTranslationRepository) DeleteProjectTranslation(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteProjectTranslation(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainProjectTranslation(q sqlc.ProjectTranslation) domain.ProjectTranslation {
	domainID := ulid.ULID(q.ID).String()
	domainIDProject := ulid.ULID(q.IDProject).String()

	var domainDescription *string
	if q.Description.Valid {
		domainDescription = &q.Description.String
	}

	return domain.ProjectTranslation{
		ID:          domainID,
		IDProject:   domainIDProject,
		Description: domainDescription,
		Name:        q.Name,
		CreatedAt:   q.CreatedAt.Time,
		UpdatedAt:   q.UpdatedAt.Time,
	}
}

func toDomainProjectTranslations(q []sqlc.ProjectTranslation) []domain.ProjectTranslation {
	domainProjectTranslations := make([]domain.ProjectTranslation, len(q))

	for i, sqlcProjectTranslation := range q {
		domainProjectTranslations[i] = toDomainProjectTranslation(sqlcProjectTranslation)
	}
	return domainProjectTranslations
}
