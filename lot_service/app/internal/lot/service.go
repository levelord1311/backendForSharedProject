package lot

import (
	"context"
	"errors"
	"fmt"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"strconv"
)

var _ Service = &service{}

type Service interface {
	Create(ctx context.Context, dto *CreateLotDTO) (uint, error)
	GetByLotID(ctx context.Context, id string) (*Lot, error)
	GetByUserID(ctx context.Context, id string) ([]*Lot, error)
	Update(ctx context.Context, dto *UpdateLotDTO) error
	Delete(ctx context.Context, lotID, userID uint) error
}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(lotStorage Storage, logger logging.Logger) (*service, error) {
	return &service{
		storage: lotStorage,
		logger:  logger,
	}, nil
}

func (s *service) Create(ctx context.Context, dto *CreateLotDTO) (uint, error) {
	lot := NewLot(dto)
	s.logger.Debug("validating lot fields...")
	if err := lot.validateFields(); err != nil {
		return 0, err
	}

	s.logger.Debug("creating new lot..")
	userID, err := s.storage.Create(ctx, lot)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to create lot. error: %w", err)
	}

	return userID, nil

}

func (s *service) GetByLotID(ctx context.Context, id string) (*Lot, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	l, err := s.storage.FindByLotID(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find lot by its id. error: %w", err)
	}
	return l, nil
}

func (s *service) GetByUserID(ctx context.Context, id string) ([]*Lot, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	l, err := s.storage.FindByUserID(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find lots by user id. error: %w", err)
	}
	return l, nil
}

func (s *service) Update(ctx context.Context, dto *UpdateLotDTO) error {
	s.logger.Debug("validating DTO fields..")
	if err := dto.validateFields(); err != nil {
		return err
	}

	updatedLot := UpdatedLot(dto)

	err := s.storage.Update(ctx, updatedLot)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update lot. error: %w", err)
	}
	return nil

}

func (s *service) Delete(ctx context.Context, lotID, userID uint) error {
	err := s.storage.Delete(ctx, lotID, userID)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete lot. error: %w", err)
	}
	return nil
}
