package domain

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	Address   *string   `json:"address,omitempty"`
	FullName  *string   `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// * Auth
type RegisterPayload struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=100"`
}

type LoginPayload struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         Role      `json:"role"`
	Address      *string   `json:"address,omitempty"`
	FullName     *string   `json:"full_name,omitempty"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateUserProfilePayload struct {
	FullName *string `json:"full_name,omitempty" form:"full_name" validate:"omitempty,min=2,max=100"`
	Address  *string `json:"address,omitempty" form:"address" validate:"omitempty,min=5,max=255"`
}

type UpdateUserPasswordPayload struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required,min=8,max=100"`
}

// * User (biasanya untuk operasi oleh Admin)
type CreateUserPayload struct {
	Username string  `json:"username" form:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string  `json:"email" form:"email" validate:"required,email"`
	Password string  `json:"password" form:"password" validate:"required,min=8,max=100"`
	Role     *Role   `json:"role,omitempty" form:"role" validate:"omitempty,oneof=admin user"`
	FullName *string `json:"full_name,omitempty" form:"full_name" validate:"omitempty,min=2,max=100"`
	Address  *string `json:"address,omitempty" form:"address" validate:"omitempty,min=5,max=255"`
}

type UpdateUserPayload struct {
	Username *string `json:"username,omitempty" form:"username" validate:"omitempty,min=3,max=30,alphanum"`
	Email    *string `json:"email,omitempty" form:"email" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" form:"password" validate:"omitempty,min=8,max=100"`
	Role     *Role   `json:"role,omitempty" form:"role" validate:"omitempty,oneof=admin user"`
	FullName *string `json:"full_name,omitempty" form:"full_name" validate:"omitempty,min=2,max=100"`
	Address  *string `json:"address,omitempty" form:"address" validate:"omitempty,min=5,max=255"`
}
