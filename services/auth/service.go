package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Rizz404/project-api-for-portfolio-website/domain"
	"github.com/Rizz404/project-api-for-portfolio-website/internal/utils"
	"github.com/oklog/ulid/v2"
)

type Repository interface {
	CreateUser(ctx context.Context, arg *domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(ctx context.Context, payload *domain.RegisterPayload) (domain.User, error) {
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
		return domain.User{}, domain.ErrConflict("Username already exist")
	}

	emailExist, err := s.repo.CheckEmailExists(ctx, params.Email)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}
	if emailExist {
		return domain.User{}, domain.ErrConflict("Email already exist")
	}

	createdUser, err := s.repo.CreateUser(ctx, &params)
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	createdUser.Password = ""
	return createdUser, nil
}

func (s *Service) Login(ctx context.Context, payload *domain.LoginPayload) (domain.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		// * Bedakan antara "user tidak ditemukan" dengan error database lainnya
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LoginResponse{}, domain.ErrUnauthorized("Invalid email or password")
		}
		// * Error lainnya adalah internal server error
		return domain.LoginResponse{}, domain.ErrInternal(err)
	}

	passwordIsValid := utils.CheckPasswordHash(payload.Password, user.Password)

	if !passwordIsValid {
		return domain.LoginResponse{}, domain.ErrUnauthorized("Invalid email or password")
	}

	jwtPayload := &utils.CreateJWTPayload{
		IDUser:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	accessToken, err := utils.CreateAccessToken(jwtPayload)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrInternal(err)
	}

	refreshToken, err := utils.CreateRefreshToken(jwtPayload.IDUser)
	if err != nil {
		return domain.LoginResponse{}, domain.ErrInternal(err)
	}

	loginResponse := &domain.LoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Address:      user.Address,
		FullName:     user.FullName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	return *loginResponse, nil
}
