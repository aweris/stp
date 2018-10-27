package sales

import (
	"context"
	"github.com/satori/go.uuid"
)

type SalesService interface {
	CreateBasket(ctx context.Context) (*uuid.UUID, error)
}
