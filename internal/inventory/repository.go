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
	AddOrUpdateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error)
	GetItemByID(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error)
	GetItemsByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.InventoryItem, error)
	FetchAllItems(ctx context.Context) ([]*models.InventoryItem, error)
}
