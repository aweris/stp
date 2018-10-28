package service

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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
		log.WithError(inventory.ErrInvalidParameter).Error("missing category")
		return nil, inventory.ErrInvalidParameter
	}

	if cat.Id != uuid.Nil {
		exist, err := is.categoryRepo.GetCategoryByID(ctx, cat.Id)
		if err != nil {
			log.WithFields(log.Fields{"category": cat}).WithError(err).Error("failed to to check category id")
			return nil, err
		}
		if exist != nil {
			log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryId).Error("category already exist with given id")
			return nil, inventory.ErrInvalidCategoryId
		}
	} else {
		cat.Id = uuid.NewV1()
	}

	if cat.Name == "" {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryName).Error("missing category name")
		return nil, inventory.ErrInvalidCategoryName
	}

	exist, err := is.categoryRepo.GetCategoryByName(ctx, cat.Name)
	if err != nil {
		log.WithFields(log.Fields{"category": cat}).WithError(err).Error("failed to check category category name")
		return nil, err
	}
	if exist != nil {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryName).Error("category already exist with given name")
		return nil, inventory.ErrInvalidCategoryName
	}

	nc, err := is.categoryRepo.SaveCategory(ctx, cat)
	if err != nil {
		log.WithFields(log.Fields{"category": cat}).WithError(err).Error("failed to add category")
		return nil, err
	}
	log.WithFields(log.Fields{"category": cat}).Info("new category added")
	return nc, nil
}

func (is *inventoryService) UpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error) {
	if cat == nil {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidParameter).Error("missing category")
		return nil, inventory.ErrInvalidParameter
	}
	if cat.Id == uuid.Nil {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryId).Error("missing category id")
		return nil, inventory.ErrInvalidCategoryId
	}
	if cat.Name == "" {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryName).Error("missing category name")
		return nil, inventory.ErrInvalidCategoryName
	}

	exist, err := is.categoryRepo.GetCategoryByID(ctx, cat.Id)
	if err != nil {
		log.WithFields(log.Fields{"category": cat}).WithError(err).Error("failed to get existing category")
		return nil, err
	}
	if exist == nil {
		log.WithFields(log.Fields{"category": cat}).WithError(inventory.ErrInvalidCategoryId).Error("failed to find category with given id")
		return nil, inventory.ErrInvalidCategoryId
	}
	nc, err := is.categoryRepo.SaveCategory(ctx, cat)
	if err != nil {
		log.WithFields(log.Fields{"category": cat}).WithError(err).Error("failed to update category")
		return nil, err
	}
	log.WithFields(log.Fields{"category": cat}).Info("category updated")
	return nc, nil
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
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(inventory.ErrInvalidCategoryId).Error("missing category id")
		return nil, inventory.ErrInvalidCategoryId
	}

	exist, err := is.categoryRepo.GetCategoryByID(ctx, categoryId)
	if err != nil {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(err).Error("failed to find category")
		return nil, err
	}
	if exist == nil {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(inventory.ErrInvalidCategoryId).Error("failed category id with given id")
		return nil, inventory.ErrInvalidCategoryId
	}

	items, err := is.itemRepo.GetItemsByCategoryID(ctx, categoryId)
	if err != nil {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(err).Error("failed to find category")
		return nil, err
	}
	if len(items) > 0 {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(inventory.ErrCategoryNotEmpty).Error("couldn't delete category")
		return nil, inventory.ErrCategoryNotEmpty
	}

	deleted, err := is.categoryRepo.DeleteCategory(ctx, categoryId)
	if err != nil {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(err).Error("failed to delete category")
		return nil, err
	}
	log.WithFields(log.Fields{"categoryId": categoryId}).Info("category deleted")
	return deleted, nil
}

func (is *inventoryService) CreateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
	if i == nil {
		log.WithError(inventory.ErrInvalidParameter).Error("missing item")
		return nil, inventory.ErrInvalidParameter
	}
	if i.Name == "" {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemName).Error("missing item name")
		return nil, inventory.ErrInvalidItemName
	}
	if i.Price.IsNegative() {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemPrice).Error("invalid item price")
		return nil, inventory.ErrInvalidItemPrice
	}

	if i.Id != uuid.Nil {
		exist, err := is.itemRepo.GetItemByID(ctx, i.Id)
		if err != nil {
			log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to check existing item with given id")
			return nil, err
		}
		if exist != nil {
			log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemId).Error("item exist with given id")
			return nil, inventory.ErrInvalidItemId
		}
	} else {
		i.Id = uuid.NewV1()
	}

	category, err := is.categoryRepo.GetCategoryByID(ctx, i.CategoryId)
	if err != nil {
		log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to item category")
		return nil, err
	}
	if category == nil {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidCategoryId).Error("missing category with given id")
		return nil, inventory.ErrInvalidCategoryId
	}

	ni, err := is.itemRepo.SaveItem(ctx, i)
	if err != nil {
		log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to create item")
		return nil, err
	}
	log.WithFields(log.Fields{"item": i}).Info("item created")
	return ni, nil
}

func (is *inventoryService) UpdateItem(ctx context.Context, i *models.InventoryItem) (*models.InventoryItem, error) {
	if i == nil {
		log.WithError(inventory.ErrInvalidParameter).Error("missing item")
		return nil, inventory.ErrInvalidParameter
	}
	if i.Id == uuid.Nil {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemId).Error("missing item id")
		return nil, inventory.ErrInvalidItemId
	}
	if i.Name == "" {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemName).Error("missing item name")
		return nil, inventory.ErrInvalidItemName
	}
	if i.Price.IsNegative() {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemPrice).Error("invalid item price")
		return nil, inventory.ErrInvalidItemPrice
	}

	exist, err := is.itemRepo.GetItemByID(ctx, i.Id)
	if err != nil {
		log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to get item")
		return nil, err
	}
	if exist == nil {
		log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidItemId).Error("failed to find item with given id")
		return nil, inventory.ErrInvalidItemId
	}

	if exist.CategoryId != i.CategoryId {
		category, err := is.categoryRepo.GetCategoryByID(ctx, i.CategoryId)
		if err != nil {
			log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to get new category")
			return nil, err
		}
		if category == nil {
			log.WithFields(log.Fields{"item": i}).WithError(inventory.ErrInvalidCategoryId).Error("failed to find new category with given id")
			return nil, inventory.ErrInvalidCategoryId
		}
	}

	ni, err := is.itemRepo.SaveItem(ctx, i)
	if err != nil {
		log.WithFields(log.Fields{"item": i}).WithError(err).Error("failed to update item")
		return nil, err
	}
	log.WithFields(log.Fields{"item": i}).Info("item updated")
	return ni, nil
}

func (is *inventoryService) GetItemByID(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error) {
	if itemId == uuid.Nil {
		log.WithFields(log.Fields{"itemId": itemId}).WithError(inventory.ErrInvalidItemId).Error("missing item id")
		return nil, inventory.ErrInvalidItemId
	}

	return is.itemRepo.GetItemByID(ctx, itemId)
}

func (is *inventoryService) GetItemsByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.InventoryItem, error) {
	if categoryId == uuid.Nil {
		log.WithFields(log.Fields{"categoryId": categoryId}).WithError(inventory.ErrInvalidCategoryId).Error("missing category id")
		return nil, inventory.ErrInvalidCategoryId
	}

	return is.itemRepo.GetItemsByCategoryID(ctx, categoryId)
}

func (is *inventoryService) FetchAllItems(ctx context.Context) ([]*models.InventoryItem, error) {
	return is.itemRepo.FetchAllItems(ctx)
}

func (is *inventoryService) DeleteItem(ctx context.Context, itemId uuid.UUID) (*models.InventoryItem, error) {
	if itemId == uuid.Nil {
		log.WithFields(log.Fields{"itemId": itemId}).WithError(inventory.ErrInvalidItemId).Error("missing item id")
		return nil, inventory.ErrInvalidItemId
	}

	return is.itemRepo.DeleteItem(ctx, itemId)
}
