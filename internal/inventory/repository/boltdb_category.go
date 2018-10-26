package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
)

const (
	bucketCategory        = "inv_category"
	bucketCategoryMeta    = "_meta"
	bucketCategoryIdx     = "index"
	bucketCategoryIdxName = "idx_category_name"
)

type boltDBCategoryRepository struct {
	db *storage.BoltDB
}

func (bcr *boltDBCategoryRepository) init() error {
	return bcr.db.Update(func(tx *bolt.Tx) error {
		cb, err := tx.CreateBucketIfNotExists([]byte(bucketCategory))
		if err != nil {
			return err
		}

		mt, err := cb.CreateBucketIfNotExists([]byte(bucketCategoryMeta))
		if err != nil {
			return err
		}

		ib, err := mt.CreateBucketIfNotExists([]byte(bucketCategoryIdx))
		if err != nil {
			return err
		}

		_, err = ib.CreateBucketIfNotExists([]byte(bucketCategoryIdxName))
		if err != nil {
			return err
		}

		return nil
	})
}

// NewBoltDBCategoryRepository creates category repository for bolt db
func NewBoltDBCategoryRepository(db *storage.BoltDB) inventory.CategoryRepository {
	bcr := &boltDBCategoryRepository{db}

	if err := bcr.init(); err != nil {
		log.Fatalln(err)
	}

	return bcr
}

// AddOrUpdateCategory adding or updating category and related indexes without checking existing value.
func (bcr *boltDBCategoryRepository) AddOrUpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error) {
	err := bcr.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		data, err := json.Marshal(cat)
		if err != nil {
			return err
		}
		err = tb.Put(cat.Id.Bytes(), data)
		if err != nil {
			return err
		}

		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		err = idx.Put([]byte(strings.ToLower(cat.Name)), cat.Id.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	return cat, err
}

// GetCategoryByID responsible for fetching category with id
func (bcr *boltDBCategoryRepository) GetCategoryByID(ctx context.Context, categoryId uuid.UUID) (*models.Category, error) {
	var t *models.Category
	err := bcr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		v := tb.Get(categoryId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &t)
		if err != nil {
			return err
		}
		return nil
	})
	return t, err
}

// GetCategoryByName responsible for fetching category with name(Name is case insensitive).
func (bcr *boltDBCategoryRepository) GetCategoryByName(ctx context.Context, categoryName string) (*models.Category, error) {
	var t *models.Category
	err := bcr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		key := idx.Get([]byte(strings.ToLower(categoryName)))
		if key == nil {
			return nil
		}

		v := tb.Get(key)
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &t)
		if err != nil {
			return err
		}
		return nil
	})
	return t, err
}

// FetchAllCategories fetching all categories
func (bcr *boltDBCategoryRepository) FetchAllCategories(ctx context.Context) ([]*models.Category, error) {
	var categories = make([]*models.Category, 0)
	err := bcr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		return tb.ForEach(func(k, v []byte) error {
			if v == nil {
				return nil
			}
			var c models.Category
			err := json.Unmarshal(v, &c)
			if err != nil {
				return err
			}
			categories = append(categories, &c)
			return nil
		})
	})
	return categories, err
}

// DeleteCategory deletes category with id
func (bcr *boltDBCategoryRepository) DeleteCategory(ctx context.Context, categoryId uuid.UUID) (*models.Category, error) {
	var existing *models.Category
	err := bcr.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		v := tb.Get(categoryId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &existing)
		if err != nil {
			return err
		}

		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		err = idx.Delete([]byte(strings.ToLower(existing.Name)))
		if err != nil {
			return err
		}

		return tb.Delete(categoryId.Bytes())
	})
	return existing, err
}
