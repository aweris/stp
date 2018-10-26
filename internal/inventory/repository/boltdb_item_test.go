package repository_test

import (
	"context"
	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"testing"
)

const (
	bucketItem                = "inv_item"
	bucketItemMeta            = "_meta"
	bucketItemIdx             = "index"
	bucketItemIdxItemCategory = "idx_item_category"
)

func TestBoltDBItemRepository_AddOrUpdateItem(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	itemId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Id:         itemId,
		Name:       "Test Item",
		CategoryId: categoryId,
		Origin:     models.ItemOriginLocal,
	}

	i, err = r.AddOrUpdateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		v := tb.Get(itemId.Bytes())
		assert.NotNil(t, v)

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketItemMeta))
		ib := mb.Bucket([]byte(bucketItemIdx))
		idx := ib.Bucket([]byte(bucketItemIdxItemCategory))

		idxIC := idx.Bucket(categoryId.Bytes())
		assert.NotNil(t, idxIC)

		idxv := idxIC.Get(i.Id.Bytes())
		assert.NotNil(t, idxv)
		return nil
	})
}

func TestBoltDBCategoryRepository_GetItemByID_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	find, err := r.GetItemByID(context.Background(), id)

	assert.NoError(t, err, "failed to find item")
	assert.Nil(t, find, "invalid item")
}

func TestBoltDBCategoryRepository_GetItemByID_ShouldReturnCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	itemId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Id:         itemId,
		Name:       "Test Item",
		CategoryId: categoryId,
		Origin:     models.ItemOriginLocal,
	}

	i, err = r.AddOrUpdateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	find, err := r.GetItemByID(context.Background(), itemId)

	assert.NoError(t, err, "failed to find item")
	assert.Equal(t, find, i, "invalid item")
}
