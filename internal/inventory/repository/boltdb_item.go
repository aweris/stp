package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	bolt "go.etcd.io/bbolt"
	"log"
)

const (
	bucketItem                = "inv_item"
	bucketItemMeta            = "_meta"
	bucketItemIdx             = "index"
	bucketItemIdxItemCategory = "idx_item_category"
)

type boltDBItemRepository struct {
	db *storage.BoltDB
}

func (bir *boltDBItemRepository) init() error {
	return bir.db.Update(func(tx *bolt.Tx) error {
		tb, err := tx.CreateBucketIfNotExists([]byte(bucketItem))
		if err != nil {
			return err
		}

		mt, err := tb.CreateBucketIfNotExists([]byte(bucketItemMeta))
		if err != nil {
			return err
		}

		ib, err := mt.CreateBucketIfNotExists([]byte(bucketItemIdx))
		if err != nil {
			return err
		}

		_, err = ib.CreateBucketIfNotExists([]byte(bucketItemIdxItemCategory))
		if err != nil {
			return err
		}

		return nil
	})
}

// NewBoltDBCategoryRepository creates item repository for bolt db
func NewBoltDBItemRepository(db *storage.BoltDB) inventory.ItemRepository {
	bir := &boltDBItemRepository{db}

	if err := bir.init(); err != nil {
		log.Fatalln(err)
	}

	return bir
}

func (bir *boltDBItemRepository) AddOrUpdateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
	err := bir.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		data, err := json.Marshal(i)
		if err != nil {
			return err
		}
		err = tb.Put(i.Id.Bytes(), data)
		if err != nil {
			return err
		}

		// getting index bucket
		mb := tb.Bucket([]byte(bucketItemMeta))
		ib := mb.Bucket([]byte(bucketItemIdx))
		idx := ib.Bucket([]byte(bucketItemIdxItemCategory))

		// creating index bucket for category
		idxIC, err := idx.CreateBucketIfNotExists(i.CategoryId.Bytes())
		if err != nil {
			return err
		}
		// adding item id under category index bucket
		err = idxIC.Put(i.Id.Bytes(), []byte("true"))
		if err != nil {
			return err
		}
		return nil
	})
	return i, err
}
