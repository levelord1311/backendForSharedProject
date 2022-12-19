package user

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/logging"
	"strconv"
)

var _ Service = &service{}

type Service interface {
	Create(ctx context.Context, dto *CreateUserDTO) (uint, error)
	SignIn(ctx context.Context, dto *SignInUserDTO) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdatePassword(ctx context.Context, dto *UpdateUserDTO) error
	Delete(ctx context.Context, id string) error
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
	user.Sanitize()

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

func (s *service) SignIn(ctx context.Context, dto *SignInUserDTO) (*User, error) {
	var u *User
	var err error
	if validation.Validate(dto.Login, is.Email) == nil {
		s.logger.Debug("received email as login, using FindByEmail method..")
		u, err = s.storage.FindByEmail(ctx, dto.Login)
	} else {
		s.logger.Debug("received username as login, using FindByUsername method..")
		u, err = s.storage.FindByUsername(ctx, dto.Login)
	}
	if err != nil {
		return nil, apperror.UnauthorizedError(err.Error())
	}
	if !u.ComparePassword(dto.Password) {
		return nil, errors.New("wrong login and/or password")
	}

	return u, nil

}

func (s *service) GetByID(ctx context.Context, id string) (*User, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	u, err := s.storage.FindByID(ctx, uint(userID))
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
	if err := dto.ValidateFields(); err != nil {
		return err
	}

	if dto.OldPassword == dto.NewPassword {
		return errors.New("new password should not match old one")
	}

	dto.Password = dto.NewPassword

	s.logger.Debugf("dto:%v", dto)

	if err := s.getUserAndComparePassword(ctx, dto); err != nil {
		return err
	}

	updatedUser := UpdatedUser(dto)

	s.logger.Debug("hashing and saving new password..")
	if err := updatedUser.EncryptPassword(); err != nil {
		return err
	}
	s.logger.Debug("deleting not hashed password..")
	updatedUser.Sanitize()

	if err := s.storage.Update(ctx, updatedUser); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil

}

func (s *service) Delete(ctx context.Context, id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if err = s.storage.Delete(ctx, uint(userID)); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return nil
}

func (s *service) getUserAndComparePassword(ctx context.Context, dto *UpdateUserDTO) error {
	userID := strconv.Itoa(int(dto.ID))
	s.logger.Debugf("id for db search:%v", userID)
	s.logger.Debug("get user by id")
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	s.logger.Debugf("user from db:%v", user)

	s.logger.Debug("compare hashed current password and old password..")
	if !user.ComparePassword(dto.OldPassword) {
		return errors.New("old password is incorrect")
	}

	return nil
}
