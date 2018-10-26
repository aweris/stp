package repository_test

import (
	"context"
	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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

func TestBoltDBItemRepository_SaveItem(t *testing.T) {
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
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = r.SaveItem(context.Background(), i)
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

func TestBoltDBItemRepository_GetItemByID_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	find, err := r.GetItemByID(context.Background(), id)

	assert.NoError(t, err, "failed to find item")
	assert.Nil(t, find, "invalid item")
}

func TestBoltDBItemRepository_GetItemByID_ShouldReturnItem(t *testing.T) {
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
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = r.SaveItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	find, err := r.GetItemByID(context.Background(), itemId)

	assert.NoError(t, err, "failed to find item")
	assert.NotNil(t, find, "missing item")
}

func TestBoltDBItemRepository_GetItemsByCategoryID_ShouldReturnItems(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	itemId1, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	itemId2, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId1, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	itemId3, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId2, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i1 := &models.InventoryItem{
		Id:         itemId1,
		Name:       "Test Item - 1",
		CategoryId: categoryId1,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i1, err = r.SaveItem(context.Background(), i1)
	assert.NoError(t, err, "failed to add item")

	i2 := &models.InventoryItem{
		Id:         itemId2,
		Name:       "Test Item - 2",
		CategoryId: categoryId1,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(10),
	}

	i2, err = r.SaveItem(context.Background(), i2)
	assert.NoError(t, err, "failed to add item")

	i3 := &models.InventoryItem{
		Id:         itemId3,
		Name:       "Test Item - 3",
		CategoryId: categoryId2,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i3, err = r.SaveItem(context.Background(), i3)
	assert.NoError(t, err, "failed to add item")

	listC1, err := r.GetItemsByCategoryID(context.Background(), categoryId1)

	assert.NotNil(t, listC1)
	assert.Equal(t, 2, len(listC1))

	listC2, err := r.GetItemsByCategoryID(context.Background(), categoryId2)

	assert.NotNil(t, listC2)
	assert.Equal(t, 1, len(listC2))
}

func TestBoltDBItemRepository_GetItemsByCategoryID_WithNonExistingCategory_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	find, err := r.GetItemsByCategoryID(context.Background(), categoryId)

	assert.NoError(t, err, "failed to find item")
	assert.Empty(t, find, "invalid item")
}

func TestBoltDBItemRepository_FetchAllItems_ShouldReturnItemList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	itemId1, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	itemId2, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId1, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	itemId3, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	categoryId2, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i1 := &models.InventoryItem{
		Id:         itemId1,
		Name:       "Test Item - 1",
		CategoryId: categoryId1,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i1, err = r.SaveItem(context.Background(), i1)
	assert.NoError(t, err, "failed to add item")

	i2 := &models.InventoryItem{
		Id:         itemId2,
		Name:       "Test Item - 2",
		CategoryId: categoryId1,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(10),
	}

	i2, err = r.SaveItem(context.Background(), i2)
	assert.NoError(t, err, "failed to add item")

	i3 := &models.InventoryItem{
		Id:         itemId3,
		Name:       "Test Item - 3",
		CategoryId: categoryId2,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i3, err = r.SaveItem(context.Background(), i3)
	assert.NoError(t, err, "failed to add item")

	list, err := r.FetchAllItems(context.Background())
	assert.NoError(t, err, "failed to fetch all items")

	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))
}
func TestBoltDBItemRepository_FetchAllItems_WithNoItem_ShouldReturnEmptyItemList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	list, err := r.FetchAllItems(context.Background())
	assert.NoError(t, err, "failed to fetch all items")

	assert.NotNil(t, list)
	assert.Empty(t, list)
}

func TestBoltDBItemRepository_DeleteItem_ShouldReturnDeletedItem(t *testing.T) {
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
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = r.SaveItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	deleted, err := r.DeleteItem(context.Background(), itemId)

	assert.NoError(t, err, "failed to delete item")
	assert.NotNil(t, deleted, "missing item")

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		v := tb.Get(itemId.Bytes())
		assert.Nil(t, v)

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketItemMeta))
		ib := mb.Bucket([]byte(bucketItemIdx))
		idx := ib.Bucket([]byte(bucketItemIdxItemCategory))

		idxIC := idx.Bucket(categoryId.Bytes())
		assert.NotNil(t, idxIC)

		idxv := idxIC.Get(i.Id.Bytes())
		assert.Nil(t, idxv)
		return nil
	})
}

func TestBoltDBItemRepository_DeleteItem_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	deleted, err := r.DeleteItem(context.Background(), id)

	assert.NoError(t, err, "failed to delete item")
	assert.Nil(t, deleted, "should be nil since we'r deleting non existing item")
}
