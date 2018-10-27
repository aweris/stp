package repository_test

import (
	"context"
	"github.com/aweris/stp/internal/models"
	basketRepository "github.com/aweris/stp/internal/sales/repository"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
	"testing"
)

const (
	bucketBasket = "sales_basket"
)

func TestBoltDBBasketRepository_SaveBasket(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := basketRepository.NewBoltDBBasketRepository(db.BoltDB)

	b := &models.Basket{
		Id:    uuid.NewV1(),
		State: models.BasketStateOpened,
	}

	b, err := r.SaveBasket(context.Background(), b)
	assert.NoError(t, err)
	assert.NotNil(t, b)

	db.BoltDB.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketBasket))
		v := tb.Get(b.Id.Bytes())
		assert.NotNil(t, v)
		return nil
	})
}
