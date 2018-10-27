package service

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type taxService struct {
	taxRepo taxes.TaxRepository
}

// NewTaxService creates inventory service with given repository interfaces
func NewTaxService(taxRepo taxes.TaxRepository) taxes.TaxService {
	return &taxService{taxRepo: taxRepo}
}

func (ts *taxService) CreateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error) {
	if tax == nil {
		return nil, taxes.ErrInvalidParameter
	}
	if tax.Name == "" {
		return nil, taxes.ErrInvalidTaxName
	}
	if !tax.Rate.IsPositive() {
		return nil, taxes.ErrInvalidTaxRate
	}

	if tax.Id != uuid.Nil {
		exist, err := ts.taxRepo.GetTaxByID(ctx, tax.Id)
		if err != nil {
			return nil, err
		}
		if exist != nil {
			return nil, taxes.ErrInvalidTaxId
		}
	} else {
		tax.Id = uuid.NewV1()
	}

	return ts.taxRepo.SaveTax(ctx, tax)
}

func (ts *taxService) UpdateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error) {
	if tax == nil {
		return nil, taxes.ErrInvalidParameter
	}
	if tax.Id == uuid.Nil {
		return nil, taxes.ErrInvalidTaxId
	}
	if tax.Name == "" {
		return nil, taxes.ErrInvalidTaxName
	}
	if !tax.Rate.IsPositive() {
		return nil, taxes.ErrInvalidTaxRate
	}

	exist, err := ts.taxRepo.GetTaxByID(ctx, tax.Id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.SaveTax(ctx, tax)
}

func (ts *taxService) GetTaxByID(ctx context.Context, taxId uuid.UUID) (*models.Tax, error) {
	if taxId == uuid.Nil {
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.GetTaxByID(ctx, taxId)
}

func (ts *taxService) FetchAllTaxes(ctx context.Context) ([]*models.Tax, error) {
	return ts.taxRepo.FetchAllTaxes(ctx)
}

func (ts *taxService) DeleteTax(ctx context.Context, taxId uuid.UUID) (*models.Tax, error) {
	if taxId == uuid.Nil {
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.DeleteTax(ctx, taxId)
}

func (ts *taxService) GetSaleItem(ctx context.Context, item *models.InventoryItem) (*models.SaleItem, error) {
	if item == nil {
		return nil, taxes.ErrInvalidParameter
	}

	taxes, err := ts.taxRepo.GetTaxesByItemOriginAndCategory(ctx, item.Origin, item.CategoryId)
	if err != nil {
		return nil, err
	}

	rate := decimal.Zero

	for _, tax := range taxes {
		rate = rate.Add(tax.Rate)
	}

	rate = rate.Div(decimal.NewFromFloat32(100))

	taxAmount := item.Price.Mul(rate).RoundCash(5)

	return &models.SaleItem{InventoryItem: item, Taxes: taxAmount, Gross: taxAmount.Add(item.Price)}, nil
}
