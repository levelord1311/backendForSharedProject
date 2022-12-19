package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) ValidateFields() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(
			requiredIf(u.EncryptedPassword == "")),
			validation.Length(6, 100)),
	)
}

type CreateUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	ID          uint   `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

type SignInUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUser(dto *CreateUserDTO) *User {
	return &User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func UpdatedUser(dto *UpdateUserDTO) *User {
	return &User{
		ID:       dto.ID,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) EncryptPassword() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func (u *User) RemoveEncryptedPassword() {
	u.EncryptedPassword = ""
}

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (dto *UpdateUserDTO) ValidateFields() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.OldPassword, validation.Required),
		validation.Field(&dto.NewPassword, validation.Required),
	)
}
