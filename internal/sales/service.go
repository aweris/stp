package sales

import (
	"context"
	"github.com/satori/go.uuid"
)

type SalesService interface {
	CreateBasket(ctx context.Context) (uuid.UUID, error)
	AddItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error)
	RemoveItem(ctx context.Context, basketId uuid.UUID, itemId uuid.UUID, itemCount int) (error)
	CancelBasket(ctx context.Context, basketId uuid.UUID) (error)
	//CloseBasket(ctx context.Context, basketId uuid.UUID) (error) // We'll add this after create receipt
}
