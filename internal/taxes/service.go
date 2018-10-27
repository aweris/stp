package taxes

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/satori/go.uuid"
)

type TaxService interface {
	CreateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
	UpdateTax(ctx context.Context, tax *models.Tax) (*models.Tax, error)
	GetTaxByID(ctx context.Context, taxId uuid.UUID) (*models.Tax, error)
	FetchAllTaxes(ctx context.Context) ([]*models.Tax, error)
	DeleteTax(ctx context.Context, taxId uuid.UUID) (*models.Tax, error)
}