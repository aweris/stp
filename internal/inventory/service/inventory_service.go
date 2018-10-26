package service

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type inventoryService struct {
	itemRepo     inventory.ItemRepository
	categoryRepo inventory.CategoryRepository
}

// NewInventoryService creates inventory service with given repository interfaces
func NewInventoryService(itemRepo inventory.ItemRepository, categoryRepo inventory.CategoryRepository) inventory.InventoryService {
	return &inventoryService{itemRepo: itemRepo, categoryRepo: categoryRepo}
}

func (is *inventoryService) CreateCategory(ctx context.Context, cat *models.Category) (*models.Category, error) {
	if cat == nil {
		return nil, inventory.ErrInvalidParameter
	}

	if cat.Id != uuid.Nil {
		exist, err := is.categoryRepo.GetCategoryByID(ctx, cat.Id)
		if err != nil {
			return nil, err
		}
		if exist != nil {
			return nil, inventory.ErrInvalidCategoryId
		}
	} else {
		id, err := uuid.NewV1()
		if err != nil {
			return nil, err
		}
		cat.Id = id
	}

	if cat.Name == "" {
		return nil, inventory.ErrInvalidCategoryName
	}

	exist, err := is.categoryRepo.GetCategoryByName(ctx, cat.Name)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, inventory.ErrInvalidCategoryName
	}

	return is.categoryRepo.SaveCategory(ctx, cat)
}

func (is *inventoryService) UpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error) {
	if cat == nil {
		return nil, inventory.ErrInvalidParameter
	}
	if cat.Id == uuid.Nil {
		return nil, inventory.ErrInvalidCategoryId
	}
	if cat.Name == "" {
		return nil, inventory.ErrInvalidCategoryName
	}

	exist, err := is.categoryRepo.GetCategoryByID(ctx, cat.Id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	return is.categoryRepo.SaveCategory(ctx, cat)
}

func (is *inventoryService) GetCategoryByID(ctx context.Context, categoryId uuid.UUID) (*models.Category, error) {
	if categoryId == uuid.Nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	return is.categoryRepo.GetCategoryByID(ctx, categoryId)
}

func (is *inventoryService) GetCategoryByName(ctx context.Context, categoryName string) (*models.Category, error) {
	if categoryName == "" {
		return nil, inventory.ErrInvalidCategoryName
	}

	return is.categoryRepo.GetCategoryByName(ctx, categoryName)
}

func (is *inventoryService) FetchAllCategories(ctx context.Context) ([]*models.Category, error) {
	return is.categoryRepo.FetchAllCategories(ctx)
}

func (is *inventoryService) DeleteCategory(ctx context.Context, categoryId uuid.UUID) (*models.Category, error) {
	if categoryId == uuid.Nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	exist, err := is.categoryRepo.GetCategoryByID(ctx, categoryId)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	items, err := is.itemRepo.GetItemsByCategoryID(ctx, categoryId)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return nil, inventory.ErrCategoryNotEmpty
	}

	return is.categoryRepo.DeleteCategory(ctx, categoryId)
}

func (is *inventoryService) CreateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
	if i == nil {
		return nil, inventory.ErrInvalidParameter
	}
	if i.Name == "" {
		return nil, inventory.ErrInvalidItemName
	}
	if i.Price.IsNegative() {
		return nil, inventory.ErrInvalidItemPrice
	}

	if i.Id != uuid.Nil {
		exist, err := is.itemRepo.GetItemByID(ctx, i.Id)
		if err != nil {
			return nil, err
		}
		if exist != nil {
			return nil, inventory.ErrInvalidItemId
		}
	} else {
		id, err := uuid.NewV1()
		if err != nil {
			return nil, err
		}
		i.Id = id
	}

	category, err := is.categoryRepo.GetCategoryByID(ctx, i.CategoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	return is.itemRepo.SaveItem(ctx, i)
}

func (is *inventoryService) UpdateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
	if i == nil {
		return nil, inventory.ErrInvalidParameter
	}
	if i.Id == uuid.Nil {
		return nil, inventory.ErrInvalidItemId
	}
	if i.Name == "" {
		return nil, inventory.ErrInvalidItemName
	}
	if i.Price.IsNegative() {
		return nil, inventory.ErrInvalidItemPrice
	}

	exist, err := is.itemRepo.GetItemByID(ctx, i.Id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, inventory.ErrInvalidItemId
	}

	if exist.CategoryId != i.CategoryId {
		category, err := is.categoryRepo.GetCategoryByID(ctx, i.CategoryId)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, inventory.ErrInvalidCategoryId
		}
	}

	return is.itemRepo.SaveItem(ctx, i)
}

func (is *inventoryService) GetItemByID(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error) {
	if itemId == uuid.Nil {
		return nil, inventory.ErrInvalidItemId
	}

	return is.itemRepo.GetItemByID(ctx, itemId)
}

func (is *inventoryService) GetItemsByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.InventoryItem, error) {
	if categoryId == uuid.Nil {
		return nil, inventory.ErrInvalidCategoryId
	}

	return is.itemRepo.GetItemsByCategoryID(ctx, categoryId)
}

func (is *inventoryService) FetchAllItems(ctx context.Context) ([]*models.InventoryItem, error) {
	return is.itemRepo.FetchAllItems(ctx)
}
