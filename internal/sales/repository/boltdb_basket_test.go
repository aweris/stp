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

func TestBoltDBBasketRepository_GetBasketByID_WhenIdNotExistInDB_ThanShouldNotReturnError(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := basketRepository.NewBoltDBBasketRepository(db.BoltDB)

	find, err := r.GetBasketByID(context.Background(), uuid.NewV1())

	assert.NoError(t, err)
	assert.Nil(t, find)
}

func TestBoltDBBasketRepository_GetBasketByID_WhenIDExistInDB_ThanShouldReturnBasket(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := basketRepository.NewBoltDBBasketRepository(db.BoltDB)

	b := &models.Basket{
		Id:    uuid.NewV1(),
		State: models.BasketStateOpened,
	}

	b, err := r.SaveBasket(context.Background(), b)
	assert.NoError(t, err)

	find, err := r.GetBasketByID(context.Background(), b.Id)

	assert.NoError(t, err)
	assert.NotNil(t, find)
}

func TestBoltDBBasketRepository_FetchAllBaskets(t *testing.T) {
	db := storage.NewTestDB()
	defer db.Close()

	r := basketRepository.NewBoltDBBasketRepository(db.BoltDB)

	b := &models.Basket{
		Id:    uuid.NewV1(),
		State: models.BasketStateOpened,
	}

	b, err := r.SaveBasket(context.Background(), b)
	assert.NoError(t, err)

	list, err := r.FetchAllBaskets(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
}
