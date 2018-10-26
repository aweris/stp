package inventory

import (
	"context"
	"github.com/aweris/stp/internal/models"
)

type CategoryRepository interface {
	AddOrUpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
}

type ItemRepository interface {
}
