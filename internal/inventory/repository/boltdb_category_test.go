package repository_test

import (
	"context"
	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
	"strings"

	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	bucketCategory        = "inv_category"
	bucketCategoryMeta    = "_meta"
	bucketCategoryIdx     = "index"
	bucketCategoryIdxName = "idx_category_name"
)

func TestBoltDBCategoryRepository_SaveCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = r.SaveCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		v := tb.Get(id.Bytes())

		assert.NotNil(t, v)

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		idxv := idx.Get([]byte(strings.ToLower(c.Name)))
		assert.Equal(t, id.Bytes(), idxv)
		return nil
	})
}

func TestBoltDBCategoryRepository_GetCategoryByID_ShouldReturnCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()

	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = r.SaveCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	find, err := r.GetCategoryByID(context.Background(), id)

	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, find, c, "invalid category")
}

func TestBoltDBCategoryRepository_GetCategoryByID_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	find, err := r.GetCategoryByID(context.Background(), id)

	assert.NoError(t, err, "failed to find category")
	assert.Nil(t, find, "invalid category")
}

func TestBoltDBCategoryRepository_GetCategoryByName_ShouldReturnCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()

	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = r.SaveCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	find, err := r.GetCategoryByName(context.Background(), "Test category")

	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, find, c, "invalid category")
}

func TestBoltDBCategoryRepository_GetCategoryByName_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	find, err := r.GetCategoryByName(context.Background(), "Non existing category")

	assert.NoError(t, err, "failed to find category")
	assert.Nil(t, find, "invalid category")
}

func TestBoltDBCategoryRepository_FetchAllCategories_ShouldReturnCategoryList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id1, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c1 := &models.Category{
		Id:   id1,
		Name: "Test Category 1",
	}

	c1, err = r.SaveCategory(context.Background(), c1)
	assert.NoError(t, err, "failed to add category")

	id2, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c2 := &models.Category{
		Id:   id2,
		Name: "Test Category 2",
	}

	c2, err = r.SaveCategory(context.Background(), c2)

	list, err := r.FetchAllCategories(context.Background())
	assert.NoError(t, err, "failed to fetch categories")

	assert.NotNil(t, list)
	assert.Equal(t, 2, len(list))
}

func TestBoltDBCategoryRepository_FetchAllCategories_WithNoCategory_ShouldReturnEmptyCategoryList(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	list, err := r.FetchAllCategories(context.Background())
	assert.NoError(t, err, "failed to fetch categories")

	assert.NotNil(t, list)
	assert.Empty(t, list)
}

func TestBoltDBCategoryRepository_DeleteCategory_ShouldReturnCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()

	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = r.SaveCategory(context.Background(), c)

	assert.NoError(t, err, "failed to add category")

	deleted, err := r.DeleteCategory(context.Background(), id)

	assert.NoError(t, err, "failed to delete category")
	assert.Equal(t, deleted, c, "invalid category deleted")

	//Check DB for deletion
	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		v := tb.Get(id.Bytes())

		assert.Nil(t, v)

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		idxv := idx.Get([]byte(strings.ToLower(c.Name)))
		assert.Nil(t, idxv)
		return nil
	})
}

func TestBoltDBCategoryRepository_DeleteCategory_WithNonExisting_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	deleted, err := r.DeleteCategory(context.Background(), id)

	assert.NoError(t, err, "failed to delete category")
	assert.Nil(t, deleted, "should be nil since we'r deleting non existing category")
}