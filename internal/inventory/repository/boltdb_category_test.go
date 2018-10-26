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

func TestBoltDBCategoryRepository_AddOrUpdateCategory(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()

	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = r.AddOrUpdateCategory(context.Background(), c)

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

	c, err = r.AddOrUpdateCategory(context.Background(), c)

	assert.NoError(t, err, "failed to add category")

	find, err := r.GetCategoryByID(context.Background(), id)

	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, find, c, "invalid category")
}

func TestBoltDBCategoryRepository_GetNoNExistingCategoryByID_ShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)

	id, err := uuid.NewV1()

	assert.NoError(t, err, "failed to generate id")

	assert.NoError(t, err, "failed to add category")

	find, err := r.GetCategoryByID(context.Background(), id)

	assert.NoError(t, err, "failed to find category")
	assert.Nil(t, find, "invalid category")
}
