package service

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
)

type salesService struct {
	basketRepo sales.BasketRepository

	invService inventory.InventoryService
	taxService taxes.TaxService
}

func NewSalesService(basketRepo sales.BasketRepository, invService inventory.InventoryService, taxService taxes.TaxService) sales.SalesService {
	return &salesService{basketRepo: basketRepo, taxService: taxService, invService: invService}
}

func (ss *salesService) CreateBasket(ctx context.Context) (*uuid.UUID, error) {
	b := &models.Basket{
		Id:    uuid.NewV1(),
		State: models.BasketStateOpened,
	}

	_, err := ss.basketRepo.SaveBasket(ctx, b)
	if err != nil {
		return nil, err
	}
	return &b.Id, nil
}
