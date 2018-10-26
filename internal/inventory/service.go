package inventory

import (
	"context"
	"github.com/aweris/stp/internal/models"
)

type InventoryService interface {
	CreateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
	UpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
}
