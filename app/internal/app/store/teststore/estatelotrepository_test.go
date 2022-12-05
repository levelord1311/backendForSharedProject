package teststore_test

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEstateLotRepository_CreateEstateLot(t *testing.T) {
	s := teststore.New()
	lot := model.TestLot(t)
	assert.NoError(t, s.EstateLot().CreateEstateLot(lot))
}
