package store

type Store interface {
	User() UserRepository
	EstateLot() EstateLotRepository
}
