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
	bucketReceipt = "sales_receipt"
)

type boltDBReceiptRepository struct {
	db *storage.BoltDB
}

func (rr *boltDBReceiptRepository) init() error {
	return rr.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketReceipt))
		if err != nil {
			return err
		}
		return nil
	})
}

func NewBoltDBReceiptRepository(db *storage.BoltDB) sales.ReceiptRepository {
	rr := &boltDBReceiptRepository{db}

	if err := rr.init(); err != nil {
		log.Fatalln(err)
	}

	return rr
}

func (rr *boltDBReceiptRepository) SaveReceipt(ctx context.Context, receipt *models.Receipt) (*models.Receipt, error) {
	err := rr.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketReceipt))

		data, err := json.Marshal(receipt)
		if err != nil {
			return err
		}

		return tb.Put(receipt.Id.Bytes(), data)
	})
	return receipt, err
}

func (rr *boltDBReceiptRepository) GetReceiptByID(ctx context.Context, receiptId uuid.UUID) (*models.Receipt, error) {
	var r *models.Receipt
	err := rr.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketReceipt))

		v := tb.Get(receiptId.Bytes())
		if v == nil {
			return nil
		}
		err := json.Unmarshal(v, &r)
		if err != nil {
			return err
		}
		return nil
	})
	return r, err
}
