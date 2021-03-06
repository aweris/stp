package service

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type salesService struct {
	basketRepo  sales.BasketRepository
	receiptRepo sales.ReceiptRepository

	invService inventory.InventoryService
	taxService taxes.TaxService
}

func NewSalesService(basketRepo sales.BasketRepository, receiptRepo sales.ReceiptRepository, invService inventory.InventoryService, taxService taxes.TaxService) sales.SalesService {
	return &salesService{basketRepo: basketRepo, receiptRepo: receiptRepo, taxService: taxService, invService: invService}
}

func (ss *salesService) CreateBasket(ctx context.Context) (uuid.UUID, error) {
	b := &models.Basket{
		Id:    uuid.NewV1(),
		Items: make(map[uuid.UUID]*models.BasketItem, 0),
		State: models.BasketStateOpened,
	}

	_, err := ss.basketRepo.SaveBasket(ctx, b)
	if err != nil {
		log.WithError(err).Error("failed to save basket")
		return uuid.Nil, err
	}
	return b.Id, nil
}

func (ss *salesService) GetBasketByID(ctx context.Context, basketId uuid.UUID) (*models.Basket, error) {
	if basketId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(sales.ErrInvalidBasketId).Error("missing basketId")
		return nil, sales.ErrInvalidBasketId
	}

	return ss.basketRepo.GetBasketByID(ctx, basketId)
}

func (ss *salesService) AddItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error) {
	if basketId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(sales.ErrInvalidBasketId).Error("missing basketId")
		return sales.ErrInvalidBasketId
	}
	if itemId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(inventory.ErrInvalidItemId).Error("missing itemId")
		return inventory.ErrInvalidItemId
	}
	if itemCount <= 0 {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(sales.ErrInvalidItemCount).Error("invalid item count")
		return sales.ErrInvalidItemCount
	}

	item, err := ss.invService.GetItemByID(ctx, itemId)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(err).Error("failed to get item with given id")
		return err
	}

	si, err := ss.taxService.GetSaleItem(ctx, item)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "item": item, "itemCount": itemCount}).WithError(err).Error("failed to get sale item")
		return err
	}

	basket, err := ss.basketRepo.GetBasketByID(ctx, basketId)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(err).Error("failed to get basket")
		return err
	}
	if basket == nil {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(sales.ErrInvalidBasketId).Error("failed to find basket with given id")
		return sales.ErrInvalidBasketId
	}

	if basket.State != models.BasketStateOpened {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(sales.ErrBasketNotOpen).Error("basket is not available")
		return sales.ErrBasketNotOpen
	}

	bi := basket.Items[si.Id]

	if bi != nil {
		bi.Count = itemCount + bi.Count
	} else {
		bi = &models.BasketItem{
			SaleItem: si,
			Count:    itemCount,
		}
	}

	basket.Items[si.Id] = bi

	ss.basketRepo.SaveBasket(ctx, basket)

	log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).Info("item added/updated in basket")

	return nil
}

func (ss *salesService) RemoveItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error) {
	if basketId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(sales.ErrInvalidBasketId).Error("missing basketId")
		return sales.ErrInvalidBasketId
	}
	if itemId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(inventory.ErrInvalidItemId).Error("missing itemId")
		return inventory.ErrInvalidItemId
	}
	if itemCount <= 0 {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(sales.ErrInvalidItemCount).Error("invalid item count")
		return sales.ErrInvalidItemCount
	}

	item, err := ss.invService.GetItemByID(ctx, itemId)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(err).Error("failed to get item with given id")
		return err
	}

	si, err := ss.taxService.GetSaleItem(ctx, item)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "item": item, "itemCount": itemCount}).WithError(err).Error("failed to get sale item")
		return err
	}

	basket, err := ss.basketRepo.GetBasketByID(ctx, basketId)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(err).Error("failed to get basket")
		return err
	}
	if basket == nil {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(sales.ErrInvalidBasketId).Error("failed to find basket with given id")
		return sales.ErrInvalidBasketId
	}

	if basket.State != models.BasketStateOpened {
		log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).WithError(sales.ErrBasketNotOpen).Error("basket is not available")
		return sales.ErrBasketNotOpen
	}

	bi := basket.Items[si.Id]

	if bi == nil {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(inventory.ErrInvalidItemId).Error("missing itemId in basket")
		return inventory.ErrInvalidItemId
	}

	if bi.Count < itemCount {
		log.WithFields(log.Fields{"basketId": basketId, "itemId": itemId, "itemCount": itemCount}).WithError(inventory.ErrInvalidItemId).Error("invalid item count for remove")
		return sales.ErrInvalidItemCount
	}

	bi.Count = bi.Count - itemCount

	if bi.Count != 0 {
		basket.Items[si.Id] = bi
	} else {
		delete(basket.Items, si.Id)
	}

	ss.basketRepo.SaveBasket(ctx, basket)

	log.WithFields(log.Fields{"basketId": basketId, "sale_item": si, "itemCount": itemCount}).Info("item removed/updated in basket")
	return nil
}

func (ss *salesService) CancelBasket(ctx context.Context, basketId uuid.UUID) (error) {
	if basketId == uuid.Nil {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(sales.ErrInvalidBasketId).Error("missing basketId")
		return sales.ErrInvalidBasketId
	}
	basket, err := ss.basketRepo.GetBasketByID(ctx, basketId)
	if err != nil {
		return err
	}
	if basket == nil {
		return sales.ErrInvalidBasketId
	}

	if basket.State != models.BasketStateOpened {
		return sales.ErrBasketNotOpen
	}

	basket.State = models.BasketStateCancelled

	_, err = ss.basketRepo.SaveBasket(ctx, basket)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(err).Error("failed to cancel basket")
		return err
	}
	log.WithFields(log.Fields{"basketId": basketId}).Info("basket cancelled")
	return nil
}

func (ss *salesService) CloseBasket(ctx context.Context, basketId uuid.UUID) (*models.Receipt, error) {
	if basketId == uuid.Nil {
		log.WithError(sales.ErrInvalidBasketId).Error("missing basketId")
		return nil, sales.ErrInvalidBasketId
	}
	basket, err := ss.basketRepo.GetBasketByID(ctx, basketId)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(err).Error("failed to get basket")
		return nil, err
	}
	if basket == nil {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(sales.ErrInvalidBasketId).Error("failed to find basket with given id")
		return nil, sales.ErrInvalidBasketId
	}

	if basket.State != models.BasketStateOpened {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(sales.ErrBasketNotOpen).Error("basket is not available")
		return nil, sales.ErrBasketNotOpen
	}

	if len(basket.Items) == 0 {
		log.WithFields(log.Fields{"basketId": basketId}).WithError(sales.ErrNotItemInBasket).Error("basket is empty")
		return nil, sales.ErrNotItemInBasket
	}

	items := make([]*models.BasketItem, 0, len(basket.Items))

	totalTax := decimal.Zero
	totalPrice := decimal.Zero
	totalGross := decimal.Zero

	for _, v := range basket.Items {
		items = append(items, v)
		totalTax = totalTax.Add(v.TotalTax())
		totalPrice = totalPrice.Add(v.TotalPrice())
		totalGross = totalGross.Add(v.TotalGross())
	}

	receipt := &models.Receipt{
		Id:         uuid.NewV1(),
		Items:      items,
		TotalTax:   totalTax,
		TotalPrice: totalPrice,
		TotalGross: totalGross,
	}

	receipt, err = ss.receiptRepo.SaveReceipt(ctx, receipt)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "receipt": receipt}).WithError(err).Error("failed to save receipt")
		return nil, err
	}

	basket.State = models.BasketStateClosed
	_, err = ss.basketRepo.SaveBasket(ctx, basket)
	if err != nil {
		log.WithFields(log.Fields{"basketId": basketId, "receipt": receipt}).WithError(err).Error("failed to close basket")
		return receipt, err
	}

	log.WithFields(log.Fields{"basketId": basketId, "receipt": receipt}).Info("basket closed")
	return receipt, nil
}

func (ss *salesService) GetReceiptByID(ctx context.Context, receiptId uuid.UUID) (*models.Receipt, error) {
	if receiptId == uuid.Nil {
		log.WithError(sales.ErrInvalidReceiptId).Error("missing receiptId")
		return nil, sales.ErrInvalidReceiptId
	}

	return ss.receiptRepo.GetReceiptByID(ctx, receiptId)
}

func (ss *salesService) FetchAllReceipts(ctx context.Context) ([]*models.Receipt, error) {
	return ss.receiptRepo.FetchAllReceipts(ctx)
}
