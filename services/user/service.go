package user

import (
	"context"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/utils"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	// * CREATE
	CreateUser(ctx context.Context, payload *domain.User) (domain.User, error)

	// * READ (MANY)
	GetUsersPaginated(ctx context.Context, limit int32, offset int32) ([]domain.User, error)
	GetUsersCursorFirst(ctx context.Context, limit int32) ([]domain.User, error)
	SearchUsers(ctx context.Context, searchTerm string) ([]domain.User, error)
	SearchUsersPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.User, error)
	SearchUsersByUsername(ctx context.Context, username string) ([]domain.User, error)
	SearchUsersByEmail(ctx context.Context, email string) ([]domain.User, error)
	SearchUsersByFullName(ctx context.Context, fullName string) ([]domain.User, error)

	// * READ (ONE & UTILITY)
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUserExists(ctx context.Context, id string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CountUsers(ctx context.Context) (int64, error)

	// * UPDATE
	UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error)
	UpdateUserProfile(ctx context.Context, payload *domain.User) (domain.User, error)

	// * DELETE
	DeleteUser(ctx context.Context, id string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// *===========================CREATE===========================*
func (s *Service) CreateUser(ctx context.Context, payload *domain.CreateUserPayload) (domain.User, error) {
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	params := domain.User{
		ID:       ulid.Make().String(),
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	usernameExist, err := s.repo.CheckUsernameExists(ctx, params.Username)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}
	if usernameExist {
		return domain.User{}, domain.ErrConflict("Username already exists")
	}

	emailExist, err := s.repo.CheckEmailExists(ctx, params.Email)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}
	if emailExist {
		return domain.User{}, domain.ErrConflict("Email already exists")
	}

	createdUser, err := s.repo.CreateUser(ctx, &params)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	createdUser.Password = ""
	return createdUser, nil
}

// *===========================READ (MANY)===========================*
func (s *Service) GetUsersPaginated(ctx context.Context, limit int32, offset int32) ([]domain.User, error) {
	users, err := s.repo.GetUsersPaginated(ctx, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) GetUsersCursorFirst(ctx context.Context, limit int32) ([]domain.User, error) {
	users, err := s.repo.GetUsersCursorFirst(ctx, limit)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) SearchUsers(ctx context.Context, searchTerm string) ([]domain.User, error) {
	users, err := s.repo.SearchUsers(ctx, searchTerm)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) SearchUsersPaginated(ctx context.Context, searchTerm string, limit int32, offset int32) ([]domain.User, error) {
	users, err := s.repo.SearchUsersPaginated(ctx, searchTerm, limit, offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) SearchUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	users, err := s.repo.SearchUsersByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) SearchUsersByEmail(ctx context.Context, email string) ([]domain.User, error) {
	users, err := s.repo.SearchUsersByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) SearchUsersByFullName(ctx context.Context, fullName string) ([]domain.User, error) {
	users, err := s.repo.SearchUsersByFullName(ctx, fullName)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// *===========================READ (ONE & UTILITY)===========================*
func (s *Service) GetUser(ctx context.Context, id string) (domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	user.Password = ""
	return user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	user.Password = ""
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	user.Password = ""
	return user, nil
}

func (s *Service) CheckUserExists(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.CheckUserExists(ctx, id)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := s.repo.CheckUsernameExists(ctx, username)
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return exists, nil
}

func (s *Service) CountUsers(ctx context.Context) (int64, error) {
	count, err := s.repo.CountUsers(ctx)
	if err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

// *===========================UPDATE===========================*
func (s *Service) UpdateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	exists, err := s.repo.CheckUserExists(ctx, payload.ID)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.User{}, domain.ErrNotFound("User")
	}

	if payload.Password != "" {
		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			return domain.User{}, domain.ErrInternal(err)
		}
		payload.Password = hashedPassword
	}

	if payload.Username != "" {
		usernameExists, err := s.repo.CheckUsernameExists(ctx, payload.Username)
		if err != nil {
			return domain.User{}, domain.ErrInternal(err)
		}
		if usernameExists {
			currentUser, err := s.repo.GetUser(ctx, payload.ID)
			if err != nil {
				return domain.User{}, domain.ErrInternal(err)
			}
			if currentUser.Username != payload.Username {
				return domain.User{}, domain.ErrConflict("Username already exists")
			}
		}
	}

	if payload.Email != "" {
		emailExists, err := s.repo.CheckEmailExists(ctx, payload.Email)
		if err != nil {
			return domain.User{}, domain.ErrInternal(err)
		}
		if emailExists {
			currentUser, err := s.repo.GetUser(ctx, payload.ID)
			if err != nil {
				return domain.User{}, domain.ErrInternal(err)
			}
			if currentUser.Email != payload.Email {
				return domain.User{}, domain.ErrConflict("Email already exists")
			}
		}
	}

	updatedUser, err := s.repo.UpdateUser(ctx, payload)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	updatedUser.Password = ""
	return updatedUser, nil
}

func (s *Service) UpdateUserProfile(ctx context.Context, payload *domain.User) (domain.User, error) {
	exists, err := s.repo.CheckUserExists(ctx, payload.ID)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}
	if !exists {
		return domain.User{}, domain.ErrNotFound("User")
	}

	updatedUser, err := s.repo.UpdateUserProfile(ctx, payload)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	updatedUser.Password = ""
	return updatedUser, nil
}

// *===========================DELETE===========================*
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	exists, err := s.repo.CheckUserExists(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}
	if !exists {
		return domain.ErrNotFound("User")
	}

	err = s.repo.DeleteUser(ctx, id)
	if err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}
