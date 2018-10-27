package sales

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type SalesService interface {
	CreateBasket(ctx context.Context) (uuid.UUID, error)
	GetBasketByID(ctx context.Context, basketId uuid.UUID) (*models.Basket, error)
	AddItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error)
	RemoveItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error)
	CancelBasket(ctx context.Context, basketId uuid.UUID) (error)
	CloseBasket(ctx context.Context, basketId uuid.UUID) (*models.Receipt, error)
	GetReceiptByID(ctx context.Context, receiptId uuid.UUID) (*models.Receipt, error)
	FetchAllReceipts(ctx context.Context) ([]*models.Receipt, error)
}
