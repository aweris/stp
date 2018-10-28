package service

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/taxes"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

var (
	twenty  = decimal.NewFromFloat32(20)
	hundred = decimal.NewFromFloat32(100)
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
		log.WithError(taxes.ErrInvalidParameter).Error("missing tax")
		return nil, taxes.ErrInvalidParameter
	}
	if tax.Name == "" {
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxName).Error("missing tax name")
		return nil, taxes.ErrInvalidTaxName
	}
	if !tax.Rate.IsPositive() {
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxRate).Error("invalid tax rate")
		return nil, taxes.ErrInvalidTaxRate
	}

	if tax.Id != uuid.Nil {
		exist, err := ts.taxRepo.GetTaxByID(ctx, tax.Id)
		if err != nil {
			log.WithFields(log.Fields{"tax": tax}).WithError(err).Error("failed to check existing taxes with given id")
			return nil, err
		}
		if exist != nil {
			log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxId).Error("tax find with given id")
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
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxId).Error("missing tax id")
		return nil, taxes.ErrInvalidTaxId
	}
	if tax.Name == "" {
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxName).Error("missing tax name")
		return nil, taxes.ErrInvalidTaxName
	}
	if !tax.Rate.IsPositive() {
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxRate).Error("missing tax rate")
		return nil, taxes.ErrInvalidTaxRate
	}

	exist, err := ts.taxRepo.GetTaxByID(ctx, tax.Id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		log.WithFields(log.Fields{"tax": tax}).WithError(taxes.ErrInvalidTaxId).Error("failed to find tax with given id")
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.SaveTax(ctx, tax)
}

func (ts *taxService) GetTaxByID(ctx context.Context, taxId uuid.UUID) (*models.Tax, error) {
	if taxId == uuid.Nil {
		log.WithError(taxes.ErrInvalidTaxId).Error("missing tax id")
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.GetTaxByID(ctx, taxId)
}

func (ts *taxService) FetchAllTaxes(ctx context.Context) ([]*models.Tax, error) {
	return ts.taxRepo.FetchAllTaxes(ctx)
}

func (ts *taxService) DeleteTax(ctx context.Context, taxId uuid.UUID) (*models.Tax, error) {
	if taxId == uuid.Nil {
		log.WithError(taxes.ErrInvalidTaxId).Error("missing tax id")
		return nil, taxes.ErrInvalidTaxId
	}

	return ts.taxRepo.DeleteTax(ctx, taxId)
}

func (ts *taxService) GetSaleItem(ctx context.Context, item *models.InventoryItem) (*models.SaleItem, error) {
	if item == nil {
		log.WithError(taxes.ErrInvalidParameter).Error("missing item")
		return nil, taxes.ErrInvalidParameter
	}

	taxes, err := ts.taxRepo.GetTaxesByItemOriginAndCategory(ctx, item.Origin, item.CategoryId)
	if err != nil {
		log.WithFields(log.Fields{"item": item}).WithError(err).Error("failed to get suitable taxes for item")
		return nil, err
	}

	rate := decimal.Zero

	for _, tax := range taxes {
		rate = rate.Add(tax.Rate)
	}

	rate = rate.Div(hundred)

	taxAmount := item.Price.Mul(rate).Mul(twenty).Ceil().Div(twenty).Round(2)

	return &models.SaleItem{InventoryItem: item, Taxes: taxAmount, Gross: taxAmount.Add(item.Price)}, nil
}
