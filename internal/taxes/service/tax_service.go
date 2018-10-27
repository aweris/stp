package service

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
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
