package taxes

import (
	"context"
	"github.com/aweris/stp/internal/models"
)

type TaxRepository interface {
	SaveTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
}