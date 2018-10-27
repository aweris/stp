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
	bucketTax = "taxes_tax"
)

func TestBoltDBTaxRepository_SaveTax(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))
		v := tb.Get(tax.Id.Bytes())
		assert.NotNil(t, v)
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
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	find, err := r.GetTaxByID(context.Background(), tax.Id)
	assert.NoError(t, err, "failed to find tax")
	assert.NotNil(t, find, "missing tax")
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemExemptToTax_ThanShouldGetEmptyList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	categoryId := uuid.NewV1()

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{categoryId: true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, categoryId)
	assert.NoError(t, err, "failed to get taxes")
	assert.Empty(t, list)
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemNotInExemptList_ThanShouldReturnResult(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, uuid.NewV1())
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 1, len(list))
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemSubjectToTax_ThanShouldReturnResult(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	categoryId := uuid.NewV1()

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.SubjectToTax,
		Categories: map[uuid.UUID]bool{categoryId: true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, categoryId)
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 1, len(list))
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemNotInSubjectList_ThanShouldReturnEmptyList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.SubjectToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, uuid.NewV1())
	assert.NoError(t, err, "failed to get taxes")
	assert.Empty(t, list)
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemOriginOk_ThanShouldReturnResult(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	categoryId := uuid.NewV1()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Sales Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginImport,
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, categoryId)
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 1, len(list))
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory_WhenItemOriginNotOk_ThanShouldResultEmptyList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	categoryId := uuid.NewV1()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Sales Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginLocal,
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, categoryId)
	assert.NoError(t, err, "failed to get taxes")
	assert.Empty(t, list)
}

func TestBoltDBTaxRepository_GetTaxesByItemOriginAndCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	localOnly := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Local 5% tax",
		Rate:   decimal.NewFromFloat32(5),
		Origin: models.TaxOriginLocal,
	}

	localOnly, err := r.SaveTax(context.Background(), localOnly)
	assert.NoError(t, err, "failed to add tax")

	importOnly := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Import 15% tax",
		Rate:   decimal.NewFromFloat32(15),
		Origin: models.TaxOriginImport,
	}

	importOnly, err = r.SaveTax(context.Background(), importOnly)
	assert.NoError(t, err, "failed to add tax")

	c1 := uuid.NewV1()
	c2 := uuid.NewV1()
	c3 := uuid.NewV1()

	subjectOnlyOne := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.SubjectToTax,
		Categories: map[uuid.UUID]bool{c1: true},
	}

	subjectOnlyOne, err = r.SaveTax(context.Background(), subjectOnlyOne)
	assert.NoError(t, err, "failed to add tax")

	exemptOnlyOne := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{c3: true},
	}

	exemptOnlyOne, err = r.SaveTax(context.Background(), exemptOnlyOne)
	assert.NoError(t, err, "failed to add tax")

	l1, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginLocal, c1)
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 3, len(l1))

	l2, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, c2)
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 2, len(l2))

	l3, err := r.GetTaxesByItemOriginAndCategory(context.Background(), models.ItemOriginImported, c3)
	assert.NoError(t, err, "failed to get taxes")
	assert.Equal(t, 1, len(l3))
}

func TestBoltDBTaxRepository_FetchAllTaxes_ThanShouldReturnItemList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	list, err := r.FetchAllTaxes(context.Background())
	assert.NoError(t, err, "failed to fetch taxes")
	assert.Equal(t, 1, len(list))
}

func TestBoltDBTaxRepository_FetchAllTaxes_WithNoItem_ThanShouldReturnEmptyItemLis(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	list, err := r.FetchAllTaxes(context.Background())
	assert.NoError(t, err, "failed to fetch taxes")
	assert.Empty(t, list)
}

func TestBoltDBTaxRepository_DeleteTax_ThanShouldDeleteTaxAndReturnDeletedItem(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	tax := &models.Tax{
		Id:         uuid.NewV1(),
		Name:       "Test Sales Tax",
		Rate:       decimal.NewFromFloat32(10),
		Origin:     models.TaxOriginAll,
		Condition:  models.ExemptToTax,
		Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
	}

	tax, err := r.SaveTax(context.Background(), tax)
	assert.NoError(t, err, "failed to add tax")

	deleted, err := r.DeleteTax(context.Background(), tax.Id)
	assert.NoError(t, err, "failed to delete tax")
	assert.NotNil(t, deleted)

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))
		v := tb.Get(tax.Id.Bytes())
		assert.Nil(t, v)
		return nil
	})
}

func TestBoltDBTaxRepository_DeleteTax_WhenDeleteNonExistingId_ThanShouldNotReturnErr(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	deleted, err := r.DeleteTax(context.Background(), uuid.NewV1())
	assert.NoError(t, err, "failed to delete tax")
	assert.Nil(t, deleted)
}
