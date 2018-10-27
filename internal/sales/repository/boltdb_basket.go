package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
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

func (br *boltDBBasketRepository) GetBasketByID(ctx context.Context, basketId uuid.UUID) (*models.Basket, error) {
	var b *models.Basket
	err := br.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketBasket))

		v := tb.Get(basketId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &b)
		if err != nil {
			return err
		}
		return nil
	})
	return b, err
}

func (br *boltDBBasketRepository) FetchAllBaskets(ctx context.Context) ([]*models.Basket, error) {
	var bs = make([]*models.Basket, 0)
	err := br.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketBasket))

		return tb.ForEach(func(k, v []byte) error {
			if v == nil {
				return nil
			}
			var b models.Basket
			err := json.Unmarshal(v, &b)
			if err != nil {
				return err
			}
			bs = append(bs, &b)
			return nil
		})
	})
	return bs, err
}
