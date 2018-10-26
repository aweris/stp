package inventory

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type CategoryRepository interface {
	AddOrUpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error)
	GetCategoryByID(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
	GetCategoryByName(ctx context.Context, categoryName string) (*models.Category, error)
	FetchAllCategories(ctx context.Context) ([]*models.Category, error)
	DeleteCategory(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
}

type ItemRepository interface {
}
