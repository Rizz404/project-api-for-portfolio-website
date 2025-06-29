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

type postgresqlProjectRepository struct {
	queries *sqlc.Queries
}

func NewProjectRepository(q *sqlc.Queries) *postgresqlProjectRepository {
	return &postgresqlProjectRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlProjectRepository) CreateProject(ctx context.Context, payload *domain.Project) (domain.Project, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Project{}, err
	}
	parsedIDUser, err := ulid.Parse(payload.IDUser)
	if err != nil {
		return domain.Project{}, err
	}
	parsedIDCategory, err := ulid.Parse(payload.IDCategory)
	if err != nil {
		return domain.Project{}, err
	}

	params := sqlc.CreateProjectParams{
		ID:           uuid.UUID(parsedID),
		IDUser:       uuid.UUID(parsedIDUser),
		IDCategory:   uuid.UUID(parsedIDCategory),
		IsDeployed:   payload.IsDeployed,
		IsMaintained: payload.IsMaintained,
	}

	if payload.LiveDemo != nil {
		params.LiveDemo = sql.NullString{String: *payload.LiveDemo, Valid: true}
	}

	if payload.SourceCode != nil {
		params.SourceCode = sql.NullString{String: *payload.SourceCode, Valid: true}
	}

	sqlcProject, err := r.queries.CreateProject(ctx, params)
	if err != nil {
		return domain.Project{}, err
	}

	return toDomainProjectFromSqlcProject(sqlcProject), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlProjectRepository) GetProjectsPaginated(ctx context.Context, limit int32, offset int32) ([]domain.Project, error) {
	sqlcProjectRows, err := r.queries.GetProjectsPaginated(ctx, sqlc.GetProjectsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectsFromPaginatedRows(sqlcProjectRows)
}

func (r *postgresqlProjectRepository) GetProjectsCursorFirst(ctx context.Context, limit int32) ([]domain.Project, error) {
	sqlcProjectRows, err := r.queries.GetProjectsCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainProjectsFromCursorFirstRows(sqlcProjectRows)
}

func (r *postgresqlProjectRepository) SearchProjectsPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.Project, error) {
	sqlcProjectRows, err := r.queries.SearchProjectsPaginated(ctx, sqlc.SearchProjectsPaginatedParams{
		Column1: searchTerm,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainProjectsFromSearchRows(sqlcProjectRows)
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlProjectRepository) GetProject(ctx context.Context, id string) (domain.Project, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.Project{}, err
	}
	sqlcProjectRow, err := r.queries.GetProject(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.Project{}, err
	}
	return toDomainProjectFromGetRow(sqlcProjectRow)
}

func (r *postgresqlProjectRepository) GetProjectByTranslatedName(ctx context.Context, name string) (domain.Project, error) {
	sqlcProjectRow, err := r.queries.GetProjectByTranslatedName(ctx, name)
	if err != nil {
		return domain.Project{}, err
	}
	return toDomainProjectFromGetByTranslatedNameRow(sqlcProjectRow)
}

func (r *postgresqlProjectRepository) CheckProjectExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckProjectExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlProjectRepository) CountProjects(ctx context.Context) (int64, error) {
	return r.queries.CountProjects(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlProjectRepository) UpdateProject(ctx context.Context, payload *domain.Project) (domain.Project, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.Project{}, err
	}

	params := sqlc.UpdateProjectParams{
		ID: uuid.UUID(parsedID),
	}

	if payload.IDCategory != "" {
		parsedIDCategory, err := ulid.Parse(payload.IDCategory)
		if err != nil {
			return domain.Project{}, err
		}
		params.IDCategory = uuid.NullUUID{UUID: uuid.UUID(parsedIDCategory), Valid: true}
	}
	params.IsDeployed = sql.NullBool{Bool: payload.IsDeployed, Valid: true}
	params.IsMaintained = sql.NullBool{Bool: payload.IsMaintained, Valid: true}
	if payload.LiveDemo != nil {
		params.LiveDemo = sql.NullString{String: *payload.LiveDemo, Valid: true}
	}
	if payload.SourceCode != nil {
		params.SourceCode = sql.NullString{String: *payload.SourceCode, Valid: true}
	}

	sqlcProject, err := r.queries.UpdateProject(ctx, params)
	if err != nil {
		return domain.Project{}, err
	}

	return toDomainProjectFromSqlcProject(sqlcProject), nil
}

// *===========================DELETE===========================*
func (r *postgresqlProjectRepository) DeleteProject(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteProject(ctx, uuid.UUID(parsedID))
}

// *=======================================HELPERS=======================================*
func toDomainProject(
	id uuid.UUID, idUser uuid.UUID, idCategory uuid.UUID, isDeployed bool, isMaintained bool,
	liveDemo sql.NullString, sourceCode sql.NullString, createdAt sql.NullTime, updatedAt sql.NullTime,
	userJSON, categoryJSON,
	translationsJSON, techsJSON, imagesJSON []byte,
) (domain.Project, error) {

	// * Unmarshal JSON-encoded data
	var user domain.User
	if err := json.Unmarshal(userJSON, &user); err != nil {
		return domain.Project{}, err
	}

	var category domain.Category
	if err := json.Unmarshal(categoryJSON, &category); err != nil {
		return domain.Project{}, err
	}

	var translations []domain.ProjectTranslation
	if err := json.Unmarshal(translationsJSON, &translations); err != nil {
		return domain.Project{}, err
	}

	var techs []domain.Tech
	if err := json.Unmarshal(techsJSON, &techs); err != nil {
		return domain.Project{}, err
	}

	var images []domain.ProjectImage
	if err := json.Unmarshal(imagesJSON, &images); err != nil {
		return domain.Project{}, err
	}

	// * Mapping ke domain.Project
	project := domain.Project{
		ID:           ulid.ULID(id).String(),
		IDUser:       ulid.ULID(idUser).String(),
		IDCategory:   ulid.ULID(idCategory).String(),
		IsDeployed:   isDeployed,
		IsMaintained: isMaintained,
		CreatedAt:    createdAt.Time,
		UpdatedAt:    updatedAt.Time,
		User:         user,
		Category:     category,
		Translations: translations,
		Techs:        techs,
		Images:       images,
	}

	if liveDemo.Valid {
		project.LiveDemo = &liveDemo.String
	}
	if sourceCode.Valid {
		project.SourceCode = &sourceCode.String
	}

	return project, nil
}

// * Helper spesifik untuk GetProjectRow
func toDomainProjectFromGetRow(row sqlc.GetProjectRow) (domain.Project, error) {
	return toDomainProject(
		row.ID, row.IDUser, row.IDCategory, row.IsDeployed, row.IsMaintained,
		row.LiveDemo, row.SourceCode, row.CreatedAt, row.UpdatedAt,
		row.User, row.Category,
		row.Translations.([]byte), row.Techs.([]byte), row.Images.([]byte),
	)
}

// * Helper spesifik untuk GetProjectByTranslatedNameRow
func toDomainProjectFromGetByTranslatedNameRow(row sqlc.GetProjectByTranslatedNameRow) (domain.Project, error) {
	return toDomainProject(
		row.ID, row.IDUser, row.IDCategory, row.IsDeployed, row.IsMaintained,
		row.LiveDemo, row.SourceCode, row.CreatedAt, row.UpdatedAt,
		row.User, row.Category,
		row.Translations.([]byte), row.Techs.([]byte), row.Images.([]byte),
	)
}

// * Helper untuk list GetProjectsPaginatedRow
func toDomainProjectsFromPaginatedRows(rows []sqlc.GetProjectsPaginatedRow) ([]domain.Project, error) {
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		project, err := toDomainProject(
			row.ID, row.IDUser, row.IDCategory, row.IsDeployed, row.IsMaintained,
			row.LiveDemo, row.SourceCode, row.CreatedAt, row.UpdatedAt,
			row.User, row.Category,
			row.Translations.([]byte), row.Techs.([]byte), row.Images.([]byte),
		)
		if err != nil {
			return nil, err
		}
		projects[i] = project
	}
	return projects, nil
}

// * Helper untuk list GetProjectsCursorFirstRow
func toDomainProjectsFromCursorFirstRows(rows []sqlc.GetProjectsCursorFirstRow) ([]domain.Project, error) {
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		project, err := toDomainProject(
			row.ID, row.IDUser, row.IDCategory, row.IsDeployed, row.IsMaintained,
			row.LiveDemo, row.SourceCode, row.CreatedAt, row.UpdatedAt,
			row.User, row.Category,
			row.Translations.([]byte), row.Techs.([]byte), row.Images.([]byte),
		)
		if err != nil {
			return nil, err
		}
		projects[i] = project
	}
	return projects, nil
}

// * Helper untuk list SearchProjectsPaginatedRow
func toDomainProjectsFromSearchRows(rows []sqlc.SearchProjectsPaginatedRow) ([]domain.Project, error) {
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		project, err := toDomainProject(
			row.ID, row.IDUser, row.IDCategory, row.IsDeployed, row.IsMaintained,
			row.LiveDemo, row.SourceCode, row.CreatedAt, row.UpdatedAt,
			row.User, row.Category,
			row.Translations.([]byte), row.Techs.([]byte), row.Images.([]byte),
		)
		if err != nil {
			return nil, err
		}
		projects[i] = project
	}
	return projects, nil
}

// * Helper khusus untuk tipe sqlc.Project yang simpel (dari Create/Update)
func toDomainProjectFromSqlcProject(q sqlc.Project) domain.Project {
	project := domain.Project{
		ID:           ulid.ULID(q.ID).String(),
		IDUser:       ulid.ULID(q.IDUser).String(),
		IDCategory:   ulid.ULID(q.IDCategory).String(),
		IsDeployed:   q.IsDeployed,
		IsMaintained: q.IsMaintained,
		CreatedAt:    q.CreatedAt.Time,
		UpdatedAt:    q.UpdatedAt.Time,
		Translations: []domain.ProjectTranslation{},
		Techs:        []domain.Tech{},
		Images:       []domain.ProjectImage{},
	}
	if q.LiveDemo.Valid {
		project.LiveDemo = &q.LiveDemo.String
	}
	if q.SourceCode.Valid {
		project.SourceCode = &q.SourceCode.String
	}
	return project
}
