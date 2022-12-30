package user_service

import (
	"time"
)

type User struct {
	ID                uint      `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password,omitempty"`
	EncryptedPassword string    `json:"-"`
	GivenName         string    `json:"given_name"`
	FamilyName        string    `json:"family_name"`
	CreatedAt         time.Time `json:"created_at"`
	RedactedAt        time.Time `json:"redacted_at"`
}

// CreateUserDTO model info
// @Description user information for registering in db. All fields are required.
type CreateUserDTO struct {
	Username string `json:"username" example:"testUser1"`
	Email    string `json:"email" example:"testUser1@mail.com"` // must be formatted as valid email address
	Password string `json:"password" example:"testPassword"`    // expected length greater than 6 symbols
}

type UpdateUserDTO struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"repeat_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

// SignInUserDTO model info
// @Description user information for authentication in db. All fields are required.
type SignInUserDTO struct {
	Login    string `json:"login"` // user's email or username
	Password string `json:"password"`
}
