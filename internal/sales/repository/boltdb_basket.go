package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/storage"
	"go.etcd.io/bbolt"
	"log"
)

const (
	bucketBasket = "sales_basket"
)

type boltDBBasketRepository struct {
	db *storage.BoltDB
}

func (br *boltDBBasketRepository) init() error {
	return br.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketBasket))
		if err != nil {
			return err
		}
		return nil
	})
}

func NewBoltDBBasketRepository(db *storage.BoltDB) sales.BasketRepository {
	br := &boltDBBasketRepository{db}

	if err := br.init(); err != nil {
		log.Fatalln(err)
	}

	return br
}

func (br *boltDBBasketRepository) SaveBasket(ctx context.Context, basket *models.Basket) (*models.Basket, error) {
	err := br.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketBasket))

		data, err := json.Marshal(basket)
		if err != nil {
			return err
		}

		return tb.Put(basket.Id.Bytes(), data)
	})
	return basket, err
}
