package storage

import (
	"context"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot"
)

type Repository interface {
	Create(ctx context.Context, lot *lot.Lot) (uint, error)
	FindByLotID(ctx context.Context, id uint) (*lot.Lot, error)
	FindByUserID(ctx context.Context, id uint) ([]*lot.Lot, error)
	FindWithFilter(ctx context.Context, options QueryOptions) ([]*lot.Lot, error)
	Update(ctx context.Context, lot *lot.Lot) error
	Delete(ctx context.Context, lotID, userID uint) error
}

type QueryOptions interface {
	GetOrderBy() string
	GetFilters() map[string][]FilterOption
}
