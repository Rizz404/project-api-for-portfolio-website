package postgresql

import (
	"context"
	"database/sql"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/repository/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type postgresqlUserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *postgresqlUserRepository {
	return &postgresqlUserRepository{
		queries: q,
	}
}

// *===========================CREATE===========================*
func (r *postgresqlUserRepository) CreateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.User{}, err
	}

	params := sqlc.CreateUserParams{
		ID:       uuid.UUID(parsedID),
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     sqlc.UserRole(payload.Role),
	}

	if payload.Address != nil {
		params.Address = sql.NullString{String: *payload.Address, Valid: true}
	}
	if payload.FullName != nil {
		params.FullName = sql.NullString{String: *payload.FullName, Valid: true}
	}

	sqlcUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(sqlcUser), nil
}

// *===========================READ (MANY)===========================*
func (r *postgresqlUserRepository) GetUsersPaginated(ctx context.Context, limit int32, offset int32) ([]domain.User, error) {
	sqlcUsers, err := r.queries.GetUsersPaginated(ctx, sqlc.GetUsersPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) GetUsersCursorFirst(ctx context.Context, limit int32) ([]domain.User, error) {
	sqlcUsers, err := r.queries.GetUsersCursorFirst(ctx, limit)
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) SearchUsers(ctx context.Context, searchTerm string) ([]domain.User, error) {
	sqlcUsers, err := r.queries.SearchUsers(ctx, sql.NullString{String: searchTerm, Valid: searchTerm != ""})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) SearchUsersPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.User, error) {
	sqlcUsers, err := r.queries.SearchUsersPaginated(ctx, sqlc.SearchUsersPaginatedParams{
		Column1: sql.NullString{String: searchTerm, Valid: searchTerm != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) SearchUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	sqlcUsers, err := r.queries.SearchUsersByUsername(ctx, sql.NullString{String: username, Valid: username != ""})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) SearchUsersByEmail(ctx context.Context, email string) ([]domain.User, error) {
	sqlcUsers, err := r.queries.SearchUsersByEmail(ctx, sql.NullString{String: email, Valid: email != ""})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

func (r *postgresqlUserRepository) SearchUsersByFullName(ctx context.Context, fullName string) ([]domain.User, error) {
	sqlcUsers, err := r.queries.SearchUsersByFullName(ctx, sql.NullString{String: fullName, Valid: fullName != ""})
	if err != nil {
		return nil, err
	}
	return toDomainUsers(sqlcUsers), nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (r *postgresqlUserRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return domain.User{}, err
	}
	sqlcUser, err := r.queries.GetUser(ctx, uuid.UUID(parsedID))
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(sqlcUser), nil
}

func (r *postgresqlUserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	sqlcUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(sqlcUser), nil
}

func (r *postgresqlUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(sqlcUser), nil
}

func (r *postgresqlUserRepository) CheckUserExists(ctx context.Context, id string) (bool, error) {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return false, err
	}
	return r.queries.CheckUserExists(ctx, uuid.UUID(parsedID))
}

func (r *postgresqlUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return r.queries.CheckEmailExists(ctx, email)
}

func (r *postgresqlUserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	return r.queries.CheckUsernameExists(ctx, username)
}

func (r *postgresqlUserRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

// *===========================UPDATE===========================*
func (r *postgresqlUserRepository) UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.User{}, err
	}

	params := sqlc.UpdateUserParams{
		ID:       uuid.UUID(parsedID),
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     sqlc.UserRole(payload.Role),
	}

	if payload.Address != nil {
		params.Address = sql.NullString{String: *payload.Address, Valid: true}
	}
	if payload.FullName != nil {
		params.FullName = sql.NullString{String: *payload.FullName, Valid: true}
	}

	sqlcUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(sqlcUser), nil
}

func (r *postgresqlUserRepository) UpdateUserProfile(ctx context.Context, payload *domain.User) (domain.User, error) {
	parsedID, err := ulid.Parse(payload.ID)
	if err != nil {
		return domain.User{}, err
	}

	params := sqlc.UpdateUserProfileParams{
		ID: uuid.UUID(parsedID),
	}

	if payload.Address != nil {
		params.Address = sql.NullString{String: *payload.Address, Valid: true}
	}
	if payload.FullName != nil {
		params.FullName = sql.NullString{String: *payload.FullName, Valid: true}
	}

	sqlcUser, err := r.queries.UpdateUserProfile(ctx, params)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(sqlcUser), nil
}

// *===========================DELETE===========================*
func (r *postgresqlUserRepository) DeleteUser(ctx context.Context, id string) error {
	parsedID, err := ulid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteUser(ctx, uuid.UUID(parsedID))
}

// *===========================HELPERS===========================*
func toDomainUser(q sqlc.User) domain.User {
	var domainFullName *string
	if q.FullName.Valid {
		domainFullName = &q.FullName.String
	}

	var domainAddress *string
	if q.Address.Valid {
		domainAddress = &q.Address.String
	}

	domainID := ulid.ULID(q.ID).String()

	return domain.User{
		ID:        domainID,
		Username:  q.Username,
		Email:     q.Email,
		Password:  q.Password,
		Role:      domain.Role(q.Role),
		FullName:  domainFullName,
		Address:   domainAddress,
		CreatedAt: q.CreatedAt.Time,
		UpdatedAt: q.UpdatedAt.Time,
	}
}

func toDomainUsers(q []sqlc.User) []domain.User {
	domainUsers := make([]domain.User, len(q))

	for i, sqlcUser := range q {
		domainUsers[i] = toDomainUser(sqlcUser)
	}
	return domainUsers
}
