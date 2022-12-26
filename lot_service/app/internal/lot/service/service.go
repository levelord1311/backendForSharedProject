package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot/storage"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/api/filter"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/api/sort"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"net/url"
	"strconv"
	"strings"
)

var _ Service = &service{}

type Service interface {
	Create(ctx context.Context, dto *lot.CreateLotDTO) (uint, error)
	GetByLotID(ctx context.Context, id string) (*lot.Lot, error)
	GetByUserID(ctx context.Context, id string) ([]*lot.Lot, error)
	GetLotsWithFilter(ctx context.Context, query url.Values) ([]*lot.Lot, error)
	Update(ctx context.Context, dto *lot.UpdateLotDTO) error
	Delete(ctx context.Context, lotID, userID uint) error
}

type service struct {
	repository storage.Repository
	logger     logging.Logger
}

func NewService(lotStorage storage.Repository, logger logging.Logger) (*service, error) {
	return &service{
		repository: lotStorage,
		logger:     logger,
	}, nil
}

func (s *service) Create(ctx context.Context, dto *lot.CreateLotDTO) (uint, error) {
	lot := lot.NewLot(dto)
	s.logger.Debug("validating lot fields...")
	if err := lot.ValidateFields(); err != nil {
		return 0, err
	}

	s.logger.Debug("creating new lot..")
	userID, err := s.repository.Create(ctx, lot)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return 0, err
		}
		return 0, fmt.Errorf("failed to create lot. error: %w", err)
	}

	return userID, nil

}

func (s *service) GetByLotID(ctx context.Context, id string) (*lot.Lot, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	l, err := s.repository.FindByLotID(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find lot by its id. error: %w", err)
	}
	return l, nil
}

func (s *service) GetByUserID(ctx context.Context, id string) ([]*lot.Lot, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	l, err := s.repository.FindByUserID(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find lots by user id. error: %w", err)
	}
	return l, nil
}

func (s *service) GetLotsWithFilter(ctx context.Context, query url.Values) ([]*lot.Lot, error) {
	var l []*lot.Lot
	var err error
	var so *sort.Options

	if options, ok := ctx.Value(sort.OptionsContextKey).(sort.Options); ok {
		so = &options
	}

	fo := getFiltersFromQuery(query)

	s.logger.Debugf("GOT FILTER OPTIONS: %v", fo)
	s.logger.Debugf("GOT SORTING OPTIONS: %v", so)

	options := storage.NewOptions(so, fo)
	s.logger.Debugf("GOT OPTIONS FOR DB: %v", options)

	l, err = s.repository.FindWithFilter(ctx, options)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find lots with filter. error: %w", err)
	}
	return l, nil
}

func (s *service) Update(ctx context.Context, dto *lot.UpdateLotDTO) error {
	s.logger.Debug("validating DTO fields..")
	if err := dto.ValidateFields(); err != nil {
		return err
	}

	updatedLot := lot.UpdatedLot(dto)

	err := s.repository.Update(ctx, updatedLot)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update lot. error: %w", err)
	}
	return nil

}

func (s *service) Delete(ctx context.Context, lotID, userID uint) error {
	err := s.repository.Delete(ctx, lotID, userID)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete lot. error: %w", err)
	}
	return nil
}

func getFiltersFromQuery(query url.Values) *filter.Options {
	fo := filter.NewOptions(make(map[string][]filter.Field))

	for fltr, values := range query {
		dataType, ok := storage.FilterDataType(fltr)
		if !ok {
			continue
		}
		for _, v := range values {
			if v == "" {
				continue
			}
			f := filter.Field{}
			f.Type = dataType
			if strings.Index(v, ":") != -1 {
				split := strings.Split(v, ":")
				if operator, ok := filter.OperatorIsAllowed(split[0]); !ok {
					f.Operator = "between"
					f.Values = split
				} else {
					f.Operator = operator
					f.Values = append(f.Values, split[1])
				}
			} else {
				f.Operator = "="
				f.Values = append(f.Values, v)
			}
			fo.Fields[fltr] = append(fo.Fields[fltr], f)
		}
	}
	return fo
}
