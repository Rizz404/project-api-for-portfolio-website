package postgresql

import (
	"context"
	"database/sql"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlLanguageRepository struct {
	queries *sqlc.Queries
}

func NewLanguageRepository(q *sqlc.Queries) *postgresqlLanguageRepository {
	return &postgresqlLanguageRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlLanguageRepository) CreateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Language{}, err
	}

	params := sqlc.CreateLanguageParams{
		ID:       uuid.UUID(parsedID),
		Name:     payload.Name,
		LangCode: payload.LangCode,
	}

	sqlcLanguage, err := r.queries.CreateLanguage(ctx, params)
	if err != nil {
		return domain.Language{}, err
	}

	return toDomainLanguage(sqlcLanguage), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlLanguageRepository) GetLanguagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.GetLanguagesPaginated(ctx, sqlc.GetLanguagesPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

func (r *postgresqlLanguageRepository) GetLanguagesCursorFirst(ctx context.Context, limit int32) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.GetLanguagesCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

func (r *postgresqlLanguageRepository) SearchLanguages(ctx context.Context, searchTerm string) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.SearchLanguages(ctx, sql.NullString{String: searchTerm, Valid: searchTerm != ""})
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

func (r *postgresqlLanguageRepository) SearchLanguagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.SearchLanguagesPaginated(ctx, sqlc.SearchLanguagesPaginatedParams{
		Column1: sql.NullString{String: searchTerm, Valid: searchTerm != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

func (r *postgresqlLanguageRepository) SearchLanguagesByName(ctx context.Context, languagename string) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.SearchLanguagesByName(ctx, sql.NullString{String: languagename, Valid: languagename != ""})
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

func (r *postgresqlLanguageRepository) SearchLanguagesByLangCode(ctx context.Context, email string) ([]domain.Language, error) {
	sqlcLanguages, err := r.queries.SearchLanguagesByLangCode(ctx, sql.NullString{String: email, Valid: email != ""})
	if err != nil {
		return nil, err
	}
	return toDomainLanguages(sqlcLanguages), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlLanguageRepository) GetLanguage(ctx context.Context, id string) (domain.Language, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.Language{}, err
	}
	sqlcLanguage, err := r.queries.GetLanguage(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.Language{}, err
	}
	return toDomainLanguage(sqlcLanguage), nil
}

func (r *postgresqlLanguageRepository) GetLanguageByName(ctx context.Context, languagename string) (domain.Language, error) {
	sqlcLanguage, err := r.queries.GetLanguageByName(ctx, languagename)
	if err != nil {
		return domain.Language{}, err
	}
	return toDomainLanguage(sqlcLanguage), nil
}

func (r *postgresqlLanguageRepository) GetLanguageByLangCode(ctx context.Context, email string) (domain.Language, error) {
	sqlcLanguage, err := r.queries.GetLanguageByLangCode(ctx, email)
	if err != nil {
		return domain.Language{}, err
	}
	return toDomainLanguage(sqlcLanguage), nil
}

func (r *postgresqlLanguageRepository) CheckLanguageExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckLanguageExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlLanguageRepository) CheckLangCodeExists(ctx context.Context, email string) (bool, error) {
	return r.queries.CheckLangCodeExists(ctx, email)
}

func (r *postgresqlLanguageRepository) CheckNameExists(ctx context.Context, languagename string) (bool, error) {
	return r.queries.CheckNameExists(ctx, languagename)
}

func (r *postgresqlLanguageRepository) CountLanguages(ctx context.Context) (int64, error) {
	return r.queries.CountLanguages(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlLanguageRepository) UpdateLanguage(ctx context.Context, payload *domain.Language) (domain.Language, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Language{}, err
	}

	params := sqlc.UpdateLanguageParams{
		ID:       uuid.UUID(parsedID),
		Name:     payload.Name,
		LangCode: payload.LangCode,
	}

	sqlcLanguage, err := r.queries.UpdateLanguage(ctx, params)
	if err != nil {
		return domain.Language{}, err
	}

	return toDomainLanguage(sqlcLanguage), nil
}

// *===========================DELETE===========================*
func (r *postgresqlLanguageRepository) DeleteLanguage(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteLanguage(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainLanguage(q sqlc.Language) domain.Language {
	domainID := ulid.ULID(q.ID).String()

	return domain.Language{
		ID:        domainID,
		Name:      q.Name,
		LangCode:  q.LangCode,
		CreatedAt: q.CreatedAt.Time,
		UpdatedAt: q.UpdatedAt.Time,
	}
}

func toDomainLanguages(q []sqlc.Language) []domain.Language {
	domainLanguages := make([]domain.Language, len(q))

	for i, sqlcLanguage := range q {
		domainLanguages[i] = toDomainLanguage(sqlcLanguage)
	}
	return domainLanguages
}
