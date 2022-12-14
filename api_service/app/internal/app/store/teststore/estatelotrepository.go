package teststore

import "backendForSharedProject/internal/app/model"

type EstateLotRepository struct {
	store *Store
	lots  map[uint]*model.EstateLot
}

func (r *EstateLotRepository) CreateEstateLot(lot *model.EstateLot) error {
	if err := lot.ValidateLotFields(); err != nil {
		return err
	}

	ID := uint(len(r.lots))
	r.lots[ID] = lot
	lot.ID = ID

	return nil
}

func (r *EstateLotRepository) GetFreshEstateLots() (*[]model.EstateLot, error) {
	//TODO implement me
	panic("implement me")
}
