package inventory

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type InventoryService interface {
	CreateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
	UpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
	GetCategoryByID(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
}
