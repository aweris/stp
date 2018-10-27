package sales

import (
	"context"
	"github.com/aweris/stp/internal/models"
)

type BasketRepository interface {
	SaveBasket(ctx context.Context, basket *models.Basket) (*models.Basket, error)
}
