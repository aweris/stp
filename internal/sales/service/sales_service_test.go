package service_test

import (
	"context"
	inventoryRepository "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"
	"github.com/aweris/stp/internal/sales"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"

	salesRepository "github.com/aweris/stp/internal/sales/repository"
	salesService "github.com/aweris/stp/internal/sales/service"
	taxRepository "github.com/aweris/stp/internal/taxes/repository"
	taxService "github.com/aweris/stp/internal/taxes/service"
	"github.com/aweris/stp/storage"
)

type mockedService struct {
	sales.SalesService

	db *storage.TestDB
}

func newMockedService() *mockedService {
	db := storage.NewTestDB()

	cr := inventoryRepository.NewBoltDBCategoryRepository(db.BoltDB)
	ir := inventoryRepository.NewBoltDBItemRepository(db.BoltDB)
	is := inventoryService.NewInventoryService(ir, cr)

	tr := taxRepository.NewBoltDBTaxRepository(db.BoltDB)
	ts := taxService.NewTaxService(tr)

	br := salesRepository.NewBoltDBBasketRepository(db.BoltDB)

	ss := salesService.NewSalesService(br, is, ts)

	return &mockedService{db: db, SalesService: ss}
}

func (ms *mockedService) Close() {
	ms.db.Close()
}

func TestSalesService_CreateBasket(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	bid, err := ts.CreateBasket(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, bid)
}
