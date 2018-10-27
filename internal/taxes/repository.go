package taxes

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type TaxRepository interface {
	SaveTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
	GetTaxByID(ctx context.Context, taxId uuid.UUID) (*models.Tax, error)
	GetTaxesByItemOriginAndCategory(ctx context.Context, origin models.ItemOrigin, categoryId uuid.UUID) ([]*models.Tax, error)
	FetchAllTaxes(ctx context.Context) ([]*models.Tax, error)
}
