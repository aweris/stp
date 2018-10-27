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
	bucketTax            = "taxes_tax"
	bucketTaxMeta        = "_meta"
	bucketTaxIdx         = "index"
	bucketTaxIdxCategory = "idx_category"
)

type boltDBTaxRepository struct {
	db *storage.BoltDB
}

func (tr *boltDBTaxRepository) init() error {
	return tr.db.Update(func(tx *bolt.Tx) error {
		tb, err := tx.CreateBucketIfNotExists([]byte(bucketTax))
		if err != nil {
			return err
		}
		mt, err := tb.CreateBucketIfNotExists([]byte(bucketTaxMeta))
		if err != nil {
			return err
		}

		ib, err := mt.CreateBucketIfNotExists([]byte(bucketTaxIdx))
		if err != nil {
			return err
		}

		_, err = ib.CreateBucketIfNotExists([]byte(bucketTaxIdxCategory))
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

		// avoid unnecessary bucket access
		if len(tax.Categories) > 0 {
			mb := tb.Bucket([]byte(bucketTaxMeta))
			ib := mb.Bucket([]byte(bucketTaxIdx))
			idx := ib.Bucket([]byte(bucketTaxIdxCategory))

			idxIC := idx.Bucket(tax.Id.Bytes())

			//Just workaround update all indexes since performance isn't problem in our case
			if idxIC != nil {
				idx.DeleteBucket(tax.Id.Bytes())
			}

			// creating index bucket for category
			idxIC, err = idx.CreateBucket(tax.Id.Bytes())
			if err != nil {
				return err
			}

			// Adding categories to index
			for _, k := range tax.Categories {
				err = idxIC.Put(k.Bytes(), []byte("true"))
				if err != nil {
					return err
				}
			}
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
