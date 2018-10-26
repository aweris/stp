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

func (bir *boltDBItemRepository) SaveItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
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

func (bir *boltDBItemRepository) GetItemByID(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error) {
	var i *models.InventoryItem
	err := bir.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		v := tb.Get(itemId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &i)
		if err != nil {
			return err
		}
		return nil
	})
	return i, err
}

func (bir *boltDBItemRepository) GetItemsByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.InventoryItem, error) {
	var items = make([]*models.InventoryItem, 0)
	err := bir.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		// getting index bucket
		mb := tb.Bucket([]byte(bucketItemMeta))
		ib := mb.Bucket([]byte(bucketItemIdx))
		idx := ib.Bucket([]byte(bucketItemIdxItemCategory))

		// creating index bucket for category
		idxIC := idx.Bucket(categoryId.Bytes())
		if idxIC == nil {
			return nil
		}

		return idxIC.ForEach(func(k, v []byte) error {
			if v == nil {
				return nil
			}
			var i models.InventoryItem
			iv := tb.Get(k)
			err := json.Unmarshal(iv, &i)
			if err != nil {
				return err
			}
			items = append(items, &i)
			return nil
		})
	})
	return items, err
}

func (bir *boltDBItemRepository) FetchAllItems(ctx context.Context) ([]*models.InventoryItem, error) {
	var items = make([]*models.InventoryItem, 0)
	err := bir.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		return tb.ForEach(func(k, v []byte) error {
			if v == nil {
				return nil
			}
			var i models.InventoryItem
			err := json.Unmarshal(v, &i)
			if err != nil {
				return err
			}
			items = append(items, &i)
			return nil
		})
	})
	return items, err
}

func (bir *boltDBItemRepository) DeleteItem(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error) {
	var existing *models.InventoryItem
	err := bir.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketItem))

		v := tb.Get(itemId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &existing)
		if err != nil {
			return err
		}

		mb := tb.Bucket([]byte(bucketItemMeta))
		ib := mb.Bucket([]byte(bucketItemIdx))
		idx := ib.Bucket([]byte(bucketItemIdxItemCategory))

		err = idx.Delete([]byte(strings.ToLower(existing.Name)))
		if err != nil {
			return err
		}

		// getting index bucket for category
		idxIC := idx.Bucket(existing.CategoryId.Bytes())
		if idxIC == nil {
			return nil
		}

		// removing item id under category index bucket
		err = idxIC.Delete(existing.Id.Bytes())
		if err != nil {
			return err
		}

		return tb.Delete(itemId.Bytes())
	})
	return existing, err
}
