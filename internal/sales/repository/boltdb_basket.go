package repository

import (
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
