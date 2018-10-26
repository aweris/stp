package service_test

import (
	"context"
	"github.com/aweris/stp/internal/inventory"
	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
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
