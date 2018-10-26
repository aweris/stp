package inventory

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type CategoryRepository interface {
	AddOrUpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
	GetCategoryByID(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
}

type ItemRepository interface {
}
