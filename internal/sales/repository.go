package sales

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type BasketRepository interface {
	SaveBasket(ctx context.Context, basket *models.Basket) (*models.Basket, error)
	GetBasketByID(ctx context.Context, basketId uuid.UUID) (*models.Basket, error)
	FetchAllBaskets(ctx context.Context) ([]*models.Basket, error)
}

type ReceiptRepository interface {
	SaveReceipt(ctx context.Context, receipt *models.Receipt) (*models.Receipt, error)
}
