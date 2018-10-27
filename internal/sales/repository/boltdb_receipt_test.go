package repository_test

import (
	"context"
	"github.com/aweris/stp/internal/models"
	salesRepository "github.com/aweris/stp/internal/sales/repository"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
	"testing"
)

const (
	bucketReceipt = "sales_receipt"
)

func TestBoltDBReceiptRepository_SaveReceipt(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := salesRepository.NewBoltDBReceiptRepository(db.BoltDB)

	re := &models.Receipt{
		Id: uuid.NewV1(),
		Items: []*models.BasketItem{
			{
				SaleItem: &models.SaleItem{
					InventoryItem: &models.InventoryItem{
						Id:         uuid.UUID{},
						Name:       "Tester",
						CategoryId: uuid.UUID{},
						Origin:     models.ItemOriginImported,
						Price:      decimal.NewFromFloat32(10),
					},
					Taxes: decimal.NewFromFloat32(1),
					Gross: decimal.NewFromFloat32(11),
				},
				Count: 1,
			},
		},
		TotalTax:   decimal.NewFromFloat32(1),
		TotalPrice: decimal.NewFromFloat32(10),
		TotalGross: decimal.NewFromFloat32(11),
	}

	re, err := r.SaveReceipt(context.Background(), re)
	assert.NoError(t, err)
	assert.NotNil(t, re)

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketReceipt))
		v := tb.Get(re.Id.Bytes())
		assert.NotNil(t, v)
		return nil
	})
}

func TestBoltDBReceiptRepository_GetReceiptByID(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := salesRepository.NewBoltDBReceiptRepository(db.BoltDB)

	re := &models.Receipt{
		Id: uuid.NewV1(),
		Items: []*models.BasketItem{
			{
				SaleItem: &models.SaleItem{
					InventoryItem: &models.InventoryItem{
						Id:         uuid.UUID{},
						Name:       "Tester",
						CategoryId: uuid.UUID{},
						Origin:     models.ItemOriginImported,
						Price:      decimal.NewFromFloat32(10),
					},
					Taxes: decimal.NewFromFloat32(1),
					Gross: decimal.NewFromFloat32(11),
				},
				Count: 1,
			},
		},
		TotalTax:   decimal.NewFromFloat32(1),
		TotalPrice: decimal.NewFromFloat32(10),
		TotalGross: decimal.NewFromFloat32(11),
	}

	re, err := r.SaveReceipt(context.Background(), re)
	assert.NoError(t, err)

	find, err := r.GetReceiptByID(context.Background(), re.Id)
	assert.NoError(t, err)
	assert.NotNil(t, find)
}

func TestBoltDBReceiptRepository_GetReceiptByID_WhenIdNotExistInDB_ThenShouldNotReturnErr(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := salesRepository.NewBoltDBReceiptRepository(db.BoltDB)

	find, err := r.GetReceiptByID(context.Background(), uuid.NewV1())
	assert.NoError(t, err)
	assert.Nil(t, find)
}

func TestBoltDBReceiptRepository_FetchAllReceipts(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := salesRepository.NewBoltDBReceiptRepository(db.BoltDB)

	re := &models.Receipt{
		Id: uuid.NewV1(),
		Items: []*models.BasketItem{
			{
				SaleItem: &models.SaleItem{
					InventoryItem: &models.InventoryItem{
						Id:         uuid.UUID{},
						Name:       "Tester",
						CategoryId: uuid.UUID{},
						Origin:     models.ItemOriginImported,
						Price:      decimal.NewFromFloat32(10),
					},
					Taxes: decimal.NewFromFloat32(1),
					Gross: decimal.NewFromFloat32(11),
				},
				Count: 1,
			},
		},
		TotalTax:   decimal.NewFromFloat32(1),
		TotalPrice: decimal.NewFromFloat32(10),
		TotalGross: decimal.NewFromFloat32(11),
	}

	re, err := r.SaveReceipt(context.Background(), re)
	assert.NoError(t, err)

	list, err := r.FetchAllReceipts(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, 1, len(list))
}
