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
