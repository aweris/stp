package service_test

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	inventoryRepository "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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

	br sales.BasketRepository
	is inventory.InventoryService
	ts taxes.TaxService
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

	return &mockedService{db: db, SalesService: ss, br: br, is: is, ts: ts}
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

func TestSalesService_AddItem_WhenBucketEmpty_ThenResultMustBeSuccess(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	c := &models.Category{
		Name: "Test Category",
	}
	c, err := ts.is.CreateCategory(ctx, c)
	assert.NoError(t, err, "failed to add category")

	item := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	item, err = ts.is.CreateItem(ctx, item)
	assert.NoError(t, err, "failed to add item")

	bid, err := ts.CreateBasket(ctx)
	assert.NoError(t, err)

	err = ts.AddItem(ctx, bid, item.Id, 1)
	assert.NoError(t, err)
}

func TestSalesService_AddItem_WhenBucketHadItemAlready_ThanShouldUpdateCount(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	c := &models.Category{
		Name: "Test Category",
	}
	c, err := ts.is.CreateCategory(ctx, c)
	assert.NoError(t, err, "failed to add category")

	item := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	item, err = ts.is.CreateItem(ctx, item)
	assert.NoError(t, err, "failed to add item")

	bid, err := ts.CreateBasket(ctx)
	assert.NoError(t, err)

	err = ts.AddItem(ctx, bid, item.Id, 1)
	assert.NoError(t, err)

	basket, err := ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 1, basket.Items[item.Id].Count)

	err = ts.AddItem(ctx, bid, item.Id, 5)
	assert.NoError(t, err)

	basket, err = ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 6, basket.Items[item.Id].Count)
}

func TestSalesService_AddItem_WhenBasketIdIsNil_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.AddItem(ctx, uuid.Nil, uuid.NewV1(), 1)
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}

func TestSalesService_AddItem_WhenItemIdIsNil_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.AddItem(ctx, uuid.NewV1(), uuid.Nil, 1)
	assert.Equal(t, inventory.ErrInvalidItemId, err)
}

func TestSalesService_AddItem_WhenItemCountIsZero_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.AddItem(ctx, uuid.NewV1(), uuid.NewV1(), 0)
	assert.Equal(t, sales.ErrInvalidItemCount, err)
}

func TestSalesService_AddItem_WhenBasketIdNotExist_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.AddItem(ctx, uuid.NewV1(), uuid.NewV1(), 2)
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}
