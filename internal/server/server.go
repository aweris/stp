package server

import (
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/storage"

	inventoryRepo "github.com/aweris/stp/internal/inventory/repository"
	inventoryService "github.com/aweris/stp/internal/inventory/service"
)

// Server is a wrapper object for internal services
type Server struct {
	db               *storage.BoltDB
	InventoryService inventory.InventoryService
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

	s := &Server{
		db:               db,
		InventoryService: is,
	}

	return s
}

func (s *Server) Close() {
	s.db.Close()
}
