package repository_test

import (
	"context"
	"github.com/aweris/stp/internal/models"
	taxRepository "github.com/aweris/stp/internal/taxes/repository"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
	"testing"
)

const (
	bucketTax            = "taxes_tax"
	bucketTaxMeta        = "_meta"
	bucketTaxIdx         = "index"
	bucketTaxIdxCategory = "idx_category"
)

func TestBoltDBTaxRepository_SaveTax(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Sales Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: []uuid.UUID{uuid.NewV1()},
		},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))

		v := tb.Get(tax.Id.Bytes())
		assert.NotNil(t, v)

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketTaxMeta))
		ib := mb.Bucket([]byte(bucketTaxIdx))
		idx := ib.Bucket([]byte(bucketTaxIdxCategory))

		idxIC := idx.Bucket(tax.Id.Bytes())
		assert.NotNil(t, idxIC)

		for _, c := range tax.Categories {
			idxv := idxIC.Get(c.Bytes())
			assert.NotNil(t, idxv)
		}

		return nil
	})
}

func TestBoltDBTaxRepository_GetTaxByID_WithNonExisting_ThanShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	find, err := r.GetTaxByID(context.Background(), uuid.NewV1())

	assert.NoError(t, err, "failed to tax")
	assert.Nil(t, find, "invalid tax")
}

func TestBoltDBTaxRepository_GetTaxByID_ThanShouldReturnItem(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Sales Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: []uuid.UUID{uuid.NewV1()},
		},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	find, err := r.GetTaxByID(context.Background(), tax.Id)
	assert.NoError(t, err, "failed to find tax")
	assert.NotNil(t, find, "missing tax")
}
