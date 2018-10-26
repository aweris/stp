package repository

import (
	"github.com/aweris/stp/internal/inventory"
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
