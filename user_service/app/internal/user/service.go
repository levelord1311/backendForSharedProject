package user

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/logging"
)

var _ Service = &service{}

type Service interface {
	Create(ctx context.Context, dto *CreateUserDTO) (uint, error)
	SignIn(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	UpdatePassword(ctx context.Context, dto *UpdateUserDTO) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (*service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

func (s *service) Create(ctx context.Context, dto *CreateUserDTO) (uint, error) {
	user := NewUser(dto)
	s.logger.Debug("validating user fields...")
	if err := user.ValidateFields(); err != nil {
		return 0, err
	}

	s.logger.Debug("generating encrypted password...")
	if err := user.EncryptPassword(); err != nil {
		return 0, err
	}

	s.logger.Debug("creating new user...")
	userID, err := s.storage.Create(ctx, user)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to create user. error: %w", err)
	}

	return userID, nil

}

func (s *service) SignIn(ctx context.Context, login, password string) (*User, error) {
	u := &User{}
	if validation.Validate(login, is.Email) == nil {
		u, err := s.storage.FindByEmail(ctx, login)
		if err != nil || !u.ComparePassword(password) {
			return nil, err
		}
	} else {
		u, err := s.storage.FindByUsername(ctx, login)
		if err != nil || !u.ComparePassword(password) {
			return nil, err
		}
	}

	return u, nil

}

func (s *service) GetByID(ctx context.Context, id uint) (*User, error) {

	u, err := s.storage.FindByID(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find user by id. error: %w", err)
	}
	return u, nil
}

func (s *service) UpdatePassword(ctx context.Context, dto *UpdateUserDTO) error {

	s.logger.Debug("validating DTO fields..")
	if err := validation.ValidateStruct(dto,
		validation.Field(&dto.OldPassword, validation.Required),
		validation.Field(&dto.NewPassword, validation.Required),
	); err != nil {
		return err
	}
	if dto.OldPassword == dto.NewPassword {
		return errors.New("new password should not match old one")
	}
	s.logger.Debug("get user by id")
	user, err := s.GetByID(ctx, dto.ID)
	if err != nil {
		return err
	}

	s.logger.Debug("compare hashed current password and old password..")
	if !user.ComparePassword(dto.OldPassword) {
		return errors.New("old password is incorrect")
	}

	dto.Password = dto.NewPassword

	updatedUser := UpdatedUser(dto)

	s.logger.Debug("hashing and saving new password..")
	if err = updatedUser.EncryptPassword(); err != nil {
		return err
	}
	s.logger.Debug("deleting not hashed password..")
	updatedUser.Sanitize()

	err = s.storage.Update(ctx, updatedUser)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil

}

func (s *service) Delete(ctx context.Context, id uint) error {
	err := s.storage.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return nil
}
