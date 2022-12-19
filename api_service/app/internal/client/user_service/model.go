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

type CreateUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"repeat_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

type SignInUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
