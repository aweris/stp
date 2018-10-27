package service_test

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	taxRepository "github.com/aweris/stp/internal/taxes/repository"
	taxService "github.com/aweris/stp/internal/taxes/service"
	"github.com/aweris/stp/storage"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedService struct {
	taxes.TaxService

	db *storage.TestDB
}

func newMockedService() *mockedService {
	db := storage.NewTestDB()

	tr := taxRepository.NewBoltDBTaxRepository(db.BoltDB)

	ts := taxService.NewTaxService(tr)

	return &mockedService{db: db, TaxService: ts}
}

func (ms *mockedService) Close() {
	ms.db.Close()
}

func TestTaxService_CreateTax_WithNilParameter_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	_, err := ts.TaxService.CreateTax(context.Background(), nil)
	assert.Equal(t, err, taxes.ErrInvalidParameter)
}

func TestTaxService_CreateTax_WithEmptyName_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxName)
}

func TestTaxService_CreateTax_WhenRateIsNegative_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Rate",
		Rate:   decimal.NewFromFloat32(-10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxRate)
}

func TestTaxService_CreateTax_WhenRateIsMissing_ThenShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Rate",
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxRate)
}

func TestTaxService_CreateTax_WhenIdIsMissing_ThanShouldGenerateAndCreateTax(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Name:   "Test Rate",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)
	assert.NotNil(t, tax)
	assert.NotEqual(t, tax.Id, uuid.Nil)
}

func TestTaxService_CreateTax_WhenNonExistingId_ThanShouldCreateTax(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Rate",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)
	assert.NotNil(t, tax)
	assert.NotEqual(t, tax.Id, uuid.Nil)
}

func TestTaxService_CreateTax_WhenExistingId_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	existing := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Test Rate",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err := ts.TaxService.CreateTax(context.Background(), existing)
	assert.NoError(t, err)

	tax := &models.Tax{
		Id:     existing.Id,
		Name:   "Test Rate",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err = ts.TaxService.CreateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxId)
}

func TestTaxService_UpdateTax_WithNilParameter_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	_, err := ts.TaxService.UpdateTax(context.Background(), nil)
	assert.Equal(t, err, taxes.ErrInvalidParameter)
}

func TestTaxService_UpdateTax_WithNilId_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Name:   "Sales Tax",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}

	_, err := ts.TaxService.UpdateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxId)
}

func TestTaxService_UpdateTax_WithEmptyName_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	tax.Name = ""

	_, err = ts.TaxService.UpdateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxName)
}

func TestTaxService_UpdateTax_WhenRateIsNegative_ThanShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	tax.Rate = decimal.NewFromFloat32(-10)

	_, err = ts.TaxService.UpdateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxRate)
}

func TestTaxService_UpdateTax_WhenIdNonExist_ThenShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	_, err := ts.TaxService.UpdateTax(context.Background(), tax)
	assert.Equal(t, err, taxes.ErrInvalidTaxId)
}

func TestTaxService_UpdateTax_ShouldUpdate(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	tax.Name = "New Name"
	tax.Rate = decimal.NewFromFloat32(50)

	tax, err = ts.TaxService.UpdateTax(context.Background(), tax)
	assert.NoError(t, err)
	assert.NotNil(t, tax)
	assert.Equal(t, tax.Name, "New Name")
}

func TestTaxService_GetTaxByID_WhenIdIsNil_ThenShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	_, err := ts.GetTaxByID(context.Background(), uuid.Nil)
	assert.Equal(t, err, taxes.ErrInvalidTaxId)
}

func TestTaxService_GetTaxByID_ShouldReturnTax(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	find, err := ts.GetTaxByID(context.Background(), tax.Id)
	assert.NoError(t, err)
	assert.NotNil(t, find)
}

func TestTaxService_FetchAllTaxes_ShouldReturnList(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	list, err := ts.TaxService.FetchAllTaxes(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))

}

func TestTaxService_DeleteTax_WhenIdIsNil_ThenShouldReturnErr(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	_, err := ts.DeleteTax(context.Background(), uuid.Nil)
	assert.Equal(t, err, taxes.ErrInvalidTaxId)
}

func TestTaxService_DeleteTax_ShouldDeleteTax(t *testing.T) {
	ts := newMockedService()
	defer ts.Close()

	tax := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Will be Updated",
		Rate:   decimal.NewFromFloat32(10),
		Origin: models.TaxOriginAll,
		TaxScope: models.TaxScope{
			Condition:  models.ExemptToTax,
			Categories: map[uuid.UUID]bool{uuid.NewV1(): true},
		},
	}
	tax, err := ts.TaxService.CreateTax(context.Background(), tax)
	assert.NoError(t, err)

	deleted, err := ts.TaxService.DeleteTax(context.Background(), tax.Id)
	assert.NoError(t, err)
	assert.NotNil(t, deleted)
}
