package server

import (
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/sales"
	"github.com/aweris/stp/internal/taxes"
	"github.com/aweris/stp/storage"

	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"

	taxRepo "github.com/aweris/stp/internal/taxes/repository"
	taxService "github.com/aweris/stp/internal/taxes/service"

	salesRepository "github.com/aweris/stp/internal/sales/repository"
	salesService "github.com/aweris/stp/internal/sales/service"
)

// Server is a wrapper object for internal services
type Server struct {
	db               *storage.BoltDB
	InventoryService inventory.InventoryService
	TaxService       taxes.TaxService
	SaleService      sales.SalesService
}

// NewServer creates and configures with boltDB storage
func NewServer(storagePath string) *Server {
	db, err := storage.NewBoltDB(storagePath)
	if err != nil {
		return nil
	}
	cr := inventoryRepo.NewBoltDBCategoryRepository(db)
	ir := inventoryRepo.NewBoltDBItemRepository(db)

	is := inventoryService.NewInventoryService(ir, cr)

	tr := taxRepo.NewBoltDBTaxRepository(db)

	ts := taxService.NewTaxService(tr)

	br := salesRepository.NewBoltDBBasketRepository(db)
	rr := salesRepository.NewBoltDBReceiptRepository(db)

	ss := salesService.NewSalesService(br, rr, is, ts)

	s := &Server{
		db:               db,
		InventoryService: is,
		TaxService:       ts,
		SaleService:      ss,
	}

	return s
}

func (s *Server) Close() {
	s.db.Close()
}
