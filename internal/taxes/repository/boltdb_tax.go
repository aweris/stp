package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"go.etcd.io/bbolt"
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

func (tr *boltDBTaxRepository) GetTaxByID(ctx context.Context, taxId uuid.UUID) (*models.Tax, error) {
	var tax *models.Tax
	err := tr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))

		v := tb.Get(taxId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &tax)
		if err != nil {
			return err
		}
		return nil
	})
	return tax, err
}

func (tr *boltDBTaxRepository) GetTaxesByItemOriginAndCategory(ctx context.Context, origin models.ItemOrigin, categoryId uuid.UUID) ([]*models.Tax, error) {
	var txs = make([]*models.Tax, 0)
	err := tr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketTax))

		tb.ForEach(func(k, v []byte) error {
			var tax *models.Tax
			err := json.Unmarshal(v, &tax)
			if err != nil {
				return err
			}

			if string(tax.Origin) == string(origin) || tax.Origin == models.TaxOriginAll {
				if (tax.Condition == models.UnknownTC) {
					txs = append(txs, tax)
					return nil
				}

				exist := tax.Categories[categoryId]

				if (tax.Condition == models.SubjectToTax && exist) || (tax.Condition == models.ExemptToTax && !exist) {
					txs = append(txs, tax)
					return nil
				}
			}

			return nil
		})

		return nil
	})
	return txs, err
}
