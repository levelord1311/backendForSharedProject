package lot

import "context"

type Storage interface {
	Create(ctx context.Context, lot *Lot) (uint, error)
	FindByLotID(ctx context.Context, id uint) (*Lot, error)
	FindByUserID(ctx context.Context, id uint) ([]*Lot, error)
	Update(ctx context.Context, lot *Lot) error
	Delete(ctx context.Context, lotID, userID uint) error
}
