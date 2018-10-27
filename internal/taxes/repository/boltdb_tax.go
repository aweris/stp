package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	"github.com/aweris/stp/storage"
	bolt "go.etcd.io/bbolt"
	"log"
)

const (
	bucketTax = "taxes_tax"
)

type boltDBTaxRepository struct {
	db *storage.BoltDB
}

func (tr *boltDBTaxRepository) init() error {
	return tr.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketTax))
		if err != nil {
			return err
		}
		return nil
	})
}

func NewBoltDBTaxRepository(db *storage.BoltDB) taxes.TaxRepository {
	tr := &boltDBTaxRepository{db}

	if err := tr.init(); err != nil {
		log.Fatalln(err)
	}

	return tr
}

func (btr *boltDBTaxRepository) SaveTax(ctx context.Context, tax *models.Tax) (*models.Tax, error) {
	err := btr.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))

		data, err := json.Marshal(tax)
		if err != nil {
			return err
		}
		return tb.Put(tax.Id.Bytes(), data)
	})
	return tax, err
}
