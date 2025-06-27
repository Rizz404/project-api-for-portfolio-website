package postgresql

import (
	"context"
	"database/sql"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlUserTranslationRepository struct {
	queries *sqlc.Queries
}

func NewUserTranslationRepository(q *sqlc.Queries) *postgresqlUserTranslationRepository {
	return &postgresqlUserTranslationRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlUserTranslationRepository) CreateUserTranslation(ctx context.Context, payload *domain.UserTranslation) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	parsedIDUser, err := ulid.Parse(payload.IDUser)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	parsedIDLanguage, err := ulid.Parse(payload.IDLanguage)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	params := sqlc.CreateUserTranslationParams{
		ID:               uuid.UUID(parsedID),
		IDUser:           uuid.UUID(parsedIDUser),
		IDLanguage:       uuid.UUID(parsedIDLanguage),
		AdditionalSkills: payload.AdditionalSkills,
		Languages:        payload.Languages,
	}

	if payload.Bio != nil {
		params.Bio = sql.NullString{String: *payload.Bio, Valid: true}
	}
	if payload.AboutMe != nil {
		params.AboutMe = sql.NullString{String: *payload.AboutMe, Valid: true}
	}
	if payload.Quote != nil {
		params.Quote = sql.NullString{String: *payload.Quote, Valid: true}
	}

	sqlcUserTranslation, err := r.queries.CreateUserTranslation(ctx, params)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	return toDomainUserTranslation(sqlcUserTranslation), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlUserTranslationRepository) GetUserTranslationsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.UserTranslation, error) {
	sqlcUserTranslations, err := r.queries.GetUserTranslationsPaginated(ctx, sqlc.GetUserTranslationsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUserTranslations(sqlcUserTranslations), nil
}

func (r *postgresqlUserTranslationRepository) GetUserTranslationsByUserIDPaginated(ctx context.Context, idUser string, limit int32, offset int32) ([]domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(idUser)
	if err != nil {
		return nil, err
	}

	sqlcUserTranslations, err := r.queries.GetUserTranslationsByUserIDPaginated(ctx, sqlc.GetUserTranslationsByUserIDPaginatedParams{
		IDUser: uuid.UUID(parsedID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUserTranslations(sqlcUserTranslations), nil
}

func (r *postgresqlUserTranslationRepository) GetUserTranslationsCursorFirst(ctx context.Context, limit int32) ([]domain.UserTranslation, error) {
	sqlcUserTranslations, err := r.queries.GetUserTranslationsCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainUserTranslations(sqlcUserTranslations), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlUserTranslationRepository) GetUserTranslation(ctx context.Context, id string) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	sqlcUserTranslation, err := r.queries.GetUserTranslation(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.UserTranslation{}, err
	}
	return toDomainUserTranslation(sqlcUserTranslation), nil
}

func (r *postgresqlUserTranslationRepository) GetUserTranslationByUserIDAndLangID(ctx context.Context, idUser string, idLanguage string) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(idUser)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	parsedIDLang, err := ulid.Parse(idLanguage)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	sqlcUserTranslations, err := r.queries.GetUserTranslationByUserIDAndLangID(ctx, sqlc.GetUserTranslationByUserIDAndLangIDParams{
		IDUser:     uuid.UUID(parsedID),
		IDLanguage: uuid.UUID(parsedIDLang),
	})
	if err != nil {
		return domain.UserTranslation{}, err
	}
	return toDomainUserTranslation(sqlcUserTranslations), nil
}

func (r *postgresqlUserTranslationRepository) GetUserTranslationByUserIDAndLangName(ctx context.Context, idUser string, langName string) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(idUser)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	sqlcUserTranslations, err := r.queries.GetUserTranslationByUserIDAndLangName(ctx, sqlc.GetUserTranslationByUserIDAndLangNameParams{
		IDUser: uuid.UUID(parsedID),
		Name:   langName,
	})
	if err != nil {
		return domain.UserTranslation{}, err
	}
	return toDomainUserTranslation(sqlcUserTranslations), nil
}

func (r *postgresqlUserTranslationRepository) GetUserTranslationByUserIDAndLangCode(ctx context.Context, idUser string, langCode string) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(idUser)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	sqlcUserTranslations, err := r.queries.GetUserTranslationByUserIDAndLangCode(ctx, sqlc.GetUserTranslationByUserIDAndLangCodeParams{
		IDUser:   uuid.UUID(parsedID),
		LangCode: langCode,
	})
	if err != nil {
		return domain.UserTranslation{}, err
	}
	return toDomainUserTranslation(sqlcUserTranslations), nil
}

func (r *postgresqlUserTranslationRepository) CheckUserTranslationExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckUserTranslationExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlUserTranslationRepository) CountUserTranslations(ctx context.Context) (int64, error) {
	return r.queries.CountUserTranslations(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlUserTranslationRepository) UpdateUserTranslation(ctx context.Context, payload *domain.UserTranslation) (domain.UserTranslation, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	parsedIDUser, err := ulid.Parse(payload.IDUser)
	if err != nil {
		return domain.UserTranslation{}, err
	}
	parsedIDLanguage, err := ulid.Parse(payload.IDLanguage)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	params := sqlc.UpdateUserTranslationParams{
		ID:               uuid.UUID(parsedID),
		IDUser:           uuid.UUID(parsedIDUser),
		IDLanguage:       uuid.UUID(parsedIDLanguage),
		AdditionalSkills: payload.AdditionalSkills,
		Languages:        payload.Languages,
	}

	if payload.Bio != nil {
		params.Bio = sql.NullString{String: *payload.Bio, Valid: true}
	}
	if payload.AboutMe != nil {
		params.AboutMe = sql.NullString{String: *payload.AboutMe, Valid: true}
	}
	if payload.Quote != nil {
		params.Quote = sql.NullString{String: *payload.Quote, Valid: true}
	}

	sqlcUserTranslation, err := r.queries.UpdateUserTranslation(ctx, params)
	if err != nil {
		return domain.UserTranslation{}, err
	}

	return toDomainUserTranslation(sqlcUserTranslation), nil
}

// *===========================DELETE===========================*
func (r *postgresqlUserTranslationRepository) DeleteUserTranslation(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserTranslation(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainUserTranslation(q sqlc.UserTranslation) domain.UserTranslation {
	domainID := ulid.ULID(q.ID).String()
	domainIDUser := ulid.ULID(q.IDUser).String()
	domainIDLanguage := ulid.ULID(q.IDLanguage).String()

	return domain.UserTranslation{
		ID:               domainID,
		IDUser:           domainIDUser,
		IDLanguage:       domainIDLanguage,
		Bio:              &q.Bio.String,
		AboutMe:          &q.AboutMe.String,
		AdditionalSkills: q.AdditionalSkills,
		Languages:        q.Languages,
		Quote:            &q.Quote.String,
		CreatedAt:        q.CreatedAt.Time,
		UpdatedAt:        q.UpdatedAt.Time,
	}
}

func toDomainUserTranslations(q []sqlc.UserTranslation) []domain.UserTranslation {
	domainUserTranslations := make([]domain.UserTranslation, len(q))

	for i, sqlcUserTranslation := range q {
		domainUserTranslations[i] = toDomainUserTranslation(sqlcUserTranslation)
	}
	return domainUserTranslations
}
