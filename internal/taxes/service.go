package taxes

import (
	"context"
	"github.com/aweris/stp/internal/models"
)

type TaxService interface {
	CreateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
	UpdateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
}