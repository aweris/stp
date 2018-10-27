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
	rr sales.ReceiptRepository

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
	rr := salesRepository.NewBoltDBReceiptRepository(db.BoltDB)

	ss := salesService.NewSalesService(br, rr, is, ts)

	return &mockedService{db: db, SalesService: ss, br: br, rr: rr, is: is, ts: ts}
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

func TestSalesService_RemoveItem_WhenItemCountBiggerThanRemove_ThanShouldRemoveItemsChangeItemCount(t *testing.T) {
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

	err = ts.AddItem(ctx, bid, item.Id, 10)
	assert.NoError(t, err)

	basket, err := ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 10, basket.Items[item.Id].Count)

	err = ts.RemoveItem(ctx, bid, item.Id, 8)
	assert.NoError(t, err)

	basket, err = ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 2, basket.Items[item.Id].Count)
}

func TestSalesService_RemoveItem_WhenItemCountEqual_ThanShouldRemoveItemFromList(t *testing.T) {
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

	err = ts.AddItem(ctx, bid, item.Id, 10)
	assert.NoError(t, err)

	basket, err := ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 10, basket.Items[item.Id].Count)

	err = ts.RemoveItem(ctx, bid, item.Id, 10)
	assert.NoError(t, err)

	basket, err = ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.Nil(t, basket.Items[item.Id])
}

func TestSalesService_RemoveItem_WhenItemCountBigger_ThanShouldReturnErr(t *testing.T) {
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

	err = ts.AddItem(ctx, bid, item.Id, 10)
	assert.NoError(t, err)

	basket, err := ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 10, basket.Items[item.Id].Count)

	err = ts.RemoveItem(ctx, bid, item.Id, 18)
	assert.Equal(t, err, sales.ErrInvalidItemCount)
}

func TestSalesService_RemoveItem_WhenBasketIdIsNil_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.RemoveItem(ctx, uuid.Nil, uuid.NewV1(), 1)
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}

func TestSalesService_RemoveItem_WhenItemIdIsNil_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.RemoveItem(ctx, uuid.NewV1(), uuid.Nil, 1)
	assert.Equal(t, inventory.ErrInvalidItemId, err)
}

func TestSalesService_RemoveItem_WhenItemCountIsZero_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.RemoveItem(ctx, uuid.NewV1(), uuid.NewV1(), 0)
	assert.Equal(t, sales.ErrInvalidItemCount, err)
}

func TestSalesService_RemoveItem_WhenBasketIdNotExist_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.RemoveItem(ctx, uuid.NewV1(), uuid.NewV1(), 2)
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}

func TestSalesService_CancelBasket_WhenBasketIdIsNil_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.CancelBasket(ctx, uuid.Nil)
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}

func TestSalesService_CancelBasket_WhenBasketIdNotExist_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	err := ts.CancelBasket(ctx, uuid.NewV1())
	assert.Equal(t, sales.ErrInvalidBasketId, err)
}

func TestSalesService_CancelBasket_WhenBasketAlreadyClosedOrCanceled_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	bid, err := ts.CreateBasket(ctx)
	assert.NoError(t, err)

	err = ts.CancelBasket(ctx, bid)
	assert.NoError(t, err)

	err = ts.CancelBasket(ctx, bid)
	assert.Equal(t, sales.ErrBasketNotOpen, err)
}

func TestSalesService_CancelBasket(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	bid, err := ts.CreateBasket(ctx)
	assert.NoError(t, err)

	err = ts.CancelBasket(ctx, bid)
	assert.NoError(t, err)
}

func TestSalesService_CloseBasket_WhenBasketHasItems_ThanShouldCloseBasketWithReceipt(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	ctx := context.Background()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
	}

	_, err := ts.ts.CreateTax(context.Background(), tax)

	c := &models.Category{
		Name: "Test Category",
	}
	c, err = ts.is.CreateCategory(ctx, c)
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

	err = ts.AddItem(ctx, bid, item.Id, 10)
	assert.NoError(t, err)

	basket, err := ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.NotNil(t, basket.Items[item.Id])
	assert.Equal(t, 10, basket.Items[item.Id].Count)

	err = ts.CloseBasket(ctx, bid)

	assert.NoError(t, err)

	basket, err = ts.br.GetBasketByID(ctx, bid)
	assert.NoError(t, err)
	assert.Equal(t, string(models.BasketStateClosed), string(basket.State))
}
