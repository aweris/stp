package service_test

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedService struct {
	inventory.InventoryService

	db *storage.TestDB
}

func newMockedService() *mockedService {
	db := storage.NewTestDB()

	cr := inventoryRepo.NewBoltDBCategoryRepository(db.BoltDB)
	ir := inventoryRepo.NewBoltDBItemRepository(db.BoltDB)

	is := inventoryService.NewInventoryService(ir, cr)

	return &mockedService{db: db, InventoryService: is}
}

func (ms *mockedService) Close() {
	ms.db.Close()
}

func TestInventoryService_CreateCategory_WithNewId_ShouldCreateOne(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")
	assert.NotNil(t, c)
}

func TestInventoryService_WithoutId_ShouldCreateCategory(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	c := &models.Category{
		Name: "Test Category",
	}

	c, err := is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")
	assert.NotNil(t, c)
	assert.NotNil(t, c.Id)
}

func TestInventoryService_CreateCategory_WithExistingId_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	existing := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	existing, err = is.CreateCategory(context.Background(), existing)
	assert.NoError(t, err, "failed to add category")

	c := &models.Category{
		Id:   id,
		Name: "Test Category Invalid Id",
	}

	c, err = is.CreateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error for id")
}

func TestInventoryService_CreateCategory_WithExistingName_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	existing := &models.Category{
		Name: "Duplicate name Test Category",
	}
	existing, err := is.CreateCategory(context.Background(), existing)
	assert.NoError(t, err, "failed to add category")

	c := &models.Category{
		Name: "duplicate Name Test Category",
	}

	c, err = is.CreateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryName, "expecting error for name")
}

func TestInventoryService_CreateCategory_WithNilObject_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.CreateCategory(context.Background(), nil)
	assert.Equal(t, err, inventory.ErrInvalidParameter, "expecting error")
}

func TestInventoryService_CreateCategory_WithEmptyName_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	c := &models.Category{}
	_, err := is.CreateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryName, "expecting error")
}

func TestInventoryService_UpdateCategory_ShouldUpdate(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	c.Name = "Updated Test Category"

	c, err = is.UpdateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to update category")
	assert.NotNil(t, c)
	assert.Equal(t, c.Name, "Updated Test Category")
}

func TestInventoryService_UpdateCategory_WithNewId_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}

	c, err = is.UpdateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error for id")
}

func TestInventoryService_UpdateCategory_WithEmptyId_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	c := &models.Category{
		Name: "Test Category",
	}

	c, err := is.UpdateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error for id")
}

func TestInventoryService_UpdateCategory_WithEmptyName_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	c.Name = ""

	c, err = is.UpdateCategory(context.Background(), c)
	assert.Equal(t, err, inventory.ErrInvalidCategoryName, "expecting error for name")
}

func TestInventoryService_UpdateCategory_WithNilObject_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.UpdateCategory(context.Background(), nil)
	assert.Equal(t, err, inventory.ErrInvalidParameter, "expecting error")
}

func TestInventoryService_GetCategoryByID_ShouldReturnCategory(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	find, err := is.GetCategoryByID(context.Background(), id)
	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, c, find)
}

func TestInventoryService_GetCategoryByID_WhenIdIsNil_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.GetCategoryByID(context.Background(), uuid.Nil)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
}

func TestInventoryService_GetCategoryByName_ShouldReturnCategory(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	c := &models.Category{
		Name: "Test Category",
	}

	c, err := is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	find, err := is.GetCategoryByName(context.Background(), c.Name)
	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, c, find)
}

func TestInventoryService_GetCategoryByName_WhenIdIsNil_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.GetCategoryByName(context.Background(), "")
	assert.Equal(t, err, inventory.ErrInvalidCategoryName, "expecting error")
}

func TestInventoryService_FetchAllCategories_ShouldReturnCategoryList(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	list, err := is.FetchAllCategories(context.Background())
	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, 1, len(list))
}

func TestInventoryService_DeleteCategory_WhenCategoryExistAndEmpty_ShouldDeleteCategory(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	deleted, err := is.DeleteCategory(context.Background(), id)

	assert.NoError(t, err, "failed to delete category")
	assert.Equal(t, deleted, c)
}

func TestInventoryService_DeleteCategory_WhenIdIsNil_ShouldReturnErr(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.DeleteCategory(context.Background(), uuid.Nil)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
}

func TestInventoryService_DeleteCategory_WhenIdIsNotExist_ShouldReturnErr(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	_, err = is.DeleteCategory(context.Background(), id)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
}

func TestInventoryService_DeleteCategory_WhenHasItems_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i1 := &models.InventoryItem{
		Name:       "Test Item - 1",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}
	_, err = is.CreateItem(context.Background(), i1)
	assert.NoError(t, err, "failed to add item")

	i2 := &models.InventoryItem{
		Name:       "Test Item - 2",
		CategoryId: c.Id,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(10),
	}
	_, err = is.CreateItem(context.Background(), i2)
	assert.NoError(t, err, "failed to add item")

	_, err = is.DeleteCategory(context.Background(), id)
	assert.Equal(t, err, inventory.ErrCategoryNotEmpty, "expecting error")
}

func TestInventoryService_CreateItem_ShouldCreateItem(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item - 1",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}
	i, err = is.CreateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")
	assert.NotNil(t, i)
}

func TestInventoryService_CreateItem_WhenItemIsNil_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	i, err := is.CreateItem(context.Background(), nil)
	assert.Equal(t, err, inventory.ErrInvalidParameter, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_CreateItem_WhenItemNameIsEmpty_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemName, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_CreateItem_WhenPriceIsNegative_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(-10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemPrice, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_CreateItem_WhenIdIsExisting_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	existing := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	existing, err = is.CreateItem(context.Background(), existing)
	assert.NoError(t, err, "failed to add item")

	i := &models.InventoryItem{
		Id:         existing.Id,
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemId, "expecting error")
	assert.Nil(t, i)

}

func TestInventoryService_CreateItem_WhenCategoryIdNotExist_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: categoryId,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_ThenShouldUpdate(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	i.Price = decimal.NewFromFloat32(15)

	i, err = is.UpdateItem(context.Background(), i)
	assert.NoError(t, err, "failed to update item")
	assert.NotNil(t, i)
	assert.True(t, i.Price.Equal(decimal.NewFromFloat32(15)))
}

func TestInventoryService_UpdateItem_WhenItemIsNil_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	i, err := is.UpdateItem(context.Background(), nil)
	assert.Equal(t, err, inventory.ErrInvalidParameter, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_WhenItemNameIsEmpty_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	itemId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Id:         itemId,
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.UpdateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemName, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_WhenPriceIsNegative_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	itemId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Id:         itemId,
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(-10),
	}

	i, err = is.UpdateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemPrice, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_WhenIdIsNil_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: categoryId,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.UpdateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemId, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_WhenItemIsNotExist_ThanShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	itemId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i := &models.InventoryItem{
		Id:         itemId,
		Name:       "Test Item - 1",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}
	i, err = is.UpdateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidItemId, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_UpdateItem_WhenNewCategoryIdNotExist_ThenShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	categoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   categoryId,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	newCategoryId, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	i.CategoryId = newCategoryId

	i, err = is.UpdateItem(context.Background(), i)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
	assert.Nil(t, i)
}

func TestInventoryService_GetItemByID_ShouldReturnItem(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	find, err := is.GetItemByID(context.Background(), i.Id)
	assert.NoError(t, err, "failed to find category")
	assert.NotNil(t, find)
}

func TestInventoryService_GetItemByID_WhenIdIsNil_ShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.GetItemByID(context.Background(), uuid.Nil)
	assert.Equal(t, err, inventory.ErrInvalidItemId, "expecting error")
}

func TestInventoryService_GetCategoryByName_ShouldReturnList(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	id, err := uuid.NewV1()
	assert.NoError(t, err, "failed to generate id")

	c := &models.Category{
		Id:   id,
		Name: "Test Category",
	}
	c, err = is.CreateCategory(context.Background(), c)
	assert.NoError(t, err, "failed to add category")

	i := &models.InventoryItem{
		Name:       "Test Item",
		CategoryId: c.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(10),
	}

	i, err = is.CreateItem(context.Background(), i)
	assert.NoError(t, err, "failed to add item")

	find, err := is.GetItemsByCategoryID(context.Background(), id)
	assert.NoError(t, err, "failed to find category")
	assert.Equal(t, 1, len(find))
}

func TestInventoryService_GetItemsByCategoryID_WhenCategoryIdIsEmpty_ThanShouldReturnError(t *testing.T) {
	is := newMockedService()
	defer is.Close()

	_, err := is.GetItemsByCategoryID(context.Background(), uuid.Nil)
	assert.Equal(t, err, inventory.ErrInvalidCategoryId, "expecting error")
}
