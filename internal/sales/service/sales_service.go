package service

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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
		return uuid.Nil, err
	}
	return b.Id, nil
}

func (ss *salesService) AddItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error) {
	if basketId == uuid.Nil {
		return sales.ErrInvalidBasketId
	}
	if itemId == uuid.Nil {
		return inventory.ErrInvalidItemId
	}
	if itemCount <= 0 {
		return sales.ErrInvalidItemCount
	}

	item, err := ss.invService.GetItemByID(ctx, itemId)
	if err != nil {
		return err
	}

	si, err := ss.taxService.GetSaleItem(ctx, item)
	if err != nil {
		return err
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
	return nil
}

func (ss *salesService) RemoveItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error) {
	if basketId == uuid.Nil {
		return sales.ErrInvalidBasketId
	}
	if itemId == uuid.Nil {
		return inventory.ErrInvalidItemId
	}
	if itemCount <= 0 {
		return sales.ErrInvalidItemCount
	}

	item, err := ss.invService.GetItemByID(ctx, itemId)
	if err != nil {
		return err
	}

	si, err := ss.taxService.GetSaleItem(ctx, item)
	if err != nil {
		return err
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

	bi := basket.Items[si.Id]

	if bi == nil {
		return inventory.ErrInvalidItemId
	}

	if bi.Count < itemCount {
		return sales.ErrInvalidItemCount
	}

	bi.Count = bi.Count - itemCount

	if bi.Count != 0 {
		basket.Items[si.Id] = bi
	} else {
		delete(basket.Items, si.Id)
	}

	ss.basketRepo.SaveBasket(ctx, basket)

	return nil
}

func (ss *salesService) CancelBasket(ctx context.Context, basketId uuid.UUID) (error) {
	if basketId == uuid.Nil {
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
	return err
}

func (ss *salesService) CloseBasket(ctx context.Context, basketId uuid.UUID) (*models.Receipt, error) {
	if basketId == uuid.Nil {
		return nil, sales.ErrInvalidBasketId
	}
	basket, err := ss.basketRepo.GetBasketByID(ctx, basketId)
	if err != nil {
		return nil, err
	}
	if basket == nil {
		return nil, sales.ErrInvalidBasketId
	}

	if basket.State != models.BasketStateOpened {
		return nil, sales.ErrBasketNotOpen
	}

	if len(basket.Items) == 0 {
		return nil, sales.ErrNotItemInBasket
	}

	items := make([]*models.BasketItem, 0, len(basket.Items))

	totalTax := decimal.Zero
	totalPrice := decimal.Zero
	totalGross := decimal.Zero

	for _, v := range basket.Items {
		items = append(items, v)
		subTax := v.Taxes.Mul(decimal.NewFromFloat32(float32(v.Count)))
		totalTax = totalTax.Add(subTax)

		subPrice := v.Price.Mul(decimal.NewFromFloat32(float32(v.Count)))
		totalPrice = totalPrice.Add(subPrice)

		subGross := v.Gross.Mul(decimal.NewFromFloat32(float32(v.Count)))
		totalGross = totalGross.Add(subGross)
	}

	receipt := &models.Receipt{
		Id:         uuid.UUID{},
		Items:      items,
		TotalTax:   totalTax,
		TotalPrice: totalPrice,
		TotalGross: totalGross,
	}

	receipt, err = ss.receiptRepo.SaveReceipt(ctx, receipt)
	if err != nil {
		return nil, err
	}

	basket.State = models.BasketStateClosed
	_, err = ss.basketRepo.SaveBasket(ctx, basket)

	return receipt, err
}
