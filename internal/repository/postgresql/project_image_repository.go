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

type postgresqlProjectImageRepository struct {
	queries *sqlc.Queries
}

func NewProjectImageRepository(q *sqlc.Queries) *postgresqlProjectImageRepository {
	return &postgresqlProjectImageRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlProjectImageRepository) CreateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.ProjectImage{}, err
	}
	parsedIDProject, err := ulid.Parse(payload.IDProject)
	if err != nil {
		return domain.ProjectImage{}, err
	}

	params := sqlc.CreateProjectImageParams{
		ID:        uuid.UUID(parsedID),
		IDProject: uuid.UUID(parsedIDProject),
		FileName:  payload.FileName,
		Url:       payload.Url,
	}

	sqlcProjectImage, err := r.queries.CreateProjectImage(ctx, params)
	if err != nil {
		return domain.ProjectImage{}, err
	}

	return toDomainProjectImage(sqlcProjectImage), nil
}

func (r *postgresqlProjectImageRepository) CreateProjectImagesBatch(ctx context.Context, payload *[]domain.ProjectImage) ([]domain.ProjectImage, error) {

	batchItems := make([]domain.ProjectImage, len(*payload))

	for _, pi := range batchItems {
		parsedID, err := ulid.Parse(pi.ID)
		if err != nil {
			return nil, err
		}
		parsedIDProject, err := ulid.Parse(pi.IDProject)
		if err != nil {
			return nil, err
		}

		batchItems = append(batchItems, domain.ProjectImage{
			ID:        uuid.UUID(parsedID).String(),
			IDProject: uuid.UUID(parsedIDProject).String(),
			FileName:  pi.FileName,
			Url:       pi.Url,
		})
	}

	jsonData, err := json.Marshal(batchItems)
	if err != nil {
		return nil, err
	}

	sqlcProjectImages, err := r.queries.CreateProjectImagesBatch(ctx, json.RawMessage(jsonData))
	if err != nil {
		return nil, err
	}

	return toDomainProjectImages(sqlcProjectImages), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlProjectImageRepository) GetProjectImagesPaginated(ctx context.Context, limit int32, offset int32) ([]domain.ProjectImage, error) {
	sqlcProjectImages, err := r.queries.GetProjectImagesPaginated(ctx, sqlc.GetProjectImagesPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectImages(sqlcProjectImages), nil
}

func (r *postgresqlProjectImageRepository) GetProjectImagesCursorFirst(ctx context.Context, limit int32) ([]domain.ProjectImage, error) {
	sqlcProjectImages, err := r.queries.GetProjectImagesCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainProjectImages(sqlcProjectImages), nil
}

func (r *postgresqlProjectImageRepository) SearchProjectImages(ctx context.Context, searchTerm string) ([]domain.ProjectImage, error) {
	sqlcProjectImages, err := r.queries.SearchProjectImages(ctx, sql.NullString{String: searchTerm, Valid: searchTerm != ""})
	if err != nil {
		return nil, err
	}
	return toDomainProjectImages(sqlcProjectImages), nil
}

func (r *postgresqlProjectImageRepository) SearchProjectImagesPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.ProjectImage, error) {
	sqlcProjectImages, err := r.queries.SearchProjectImagesPaginated(ctx, sqlc.SearchProjectImagesPaginatedParams{
		Column1: sql.NullString{String: searchTerm, Valid: searchTerm != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectImages(sqlcProjectImages), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlProjectImageRepository) GetProjectImage(ctx context.Context, id string) (domain.ProjectImage, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.ProjectImage{}, err
	}
	sqlcProjectImage, err := r.queries.GetProjectImage(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.ProjectImage{}, err
	}
	return toDomainProjectImage(sqlcProjectImage), nil
}

func (r *postgresqlProjectImageRepository) CheckProjectImageExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckProjectImageExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlProjectImageRepository) CountProjectImages(ctx context.Context) (int64, error) {
	return r.queries.CountProjectImages(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlProjectImageRepository) UpdateProjectImage(ctx context.Context, payload *domain.ProjectImage) (domain.ProjectImage, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.ProjectImage{}, err
	}

	params := sqlc.UpdateProjectImageParams{
		ID:       uuid.UUID(parsedID),
		FileName: payload.FileName,
		Url:      payload.Url,
	}

	sqlcProjectImage, err := r.queries.UpdateProjectImage(ctx, params)
	if err != nil {
		return domain.ProjectImage{}, err
	}

	return toDomainProjectImage(sqlcProjectImage), nil
}

// *===========================DELETE===========================*
func (r *postgresqlProjectImageRepository) DeleteProjectImage(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteProjectImage(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainProjectImage(q sqlc.ProjectImage) domain.ProjectImage {
	domainID := ulid.ULID(q.ID).String()
	domainIDProject := ulid.ULID(q.IDProject).String()

	return domain.ProjectImage{
		ID:        domainID,
		IDProject: domainIDProject,
		FileName:  q.FileName,
		Url:       q.Url,
		CreatedAt: q.CreatedAt.Time,
		UpdatedAt: q.UpdatedAt.Time,
	}
}

func toDomainProjectImages(q []sqlc.ProjectImage) []domain.ProjectImage {
	domainProjectImages := make([]domain.ProjectImage, len(q))

	for i, sqlcProjectImage := range q {
		domainProjectImages[i] = toDomainProjectImage(sqlcProjectImage)
	}
	return domainProjectImages
}
