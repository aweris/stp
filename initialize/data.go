package initialize

import (
	"context"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/internal/server"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

func LoadTestData(server *server.Server) {

	foodC := &models.Category{
		Id:   uuid.NewV1(),
		Name: "Food",
	}
	server.InventoryService.CreateCategory(context.Background(), foodC)

	booksC := &models.Category{
		Id:   uuid.NewV1(),
		Name: "Books",
	}
	server.InventoryService.CreateCategory(context.Background(), booksC)

	musicC := &models.Category{
		Id:   uuid.NewV1(),
		Name: "Music",
	}
	server.InventoryService.CreateCategory(context.Background(), musicC)

	medicalC := &models.Category{
		Id:   uuid.NewV1(),
		Name: "Medical",
	}
	server.InventoryService.CreateCategory(context.Background(), medicalC)

	cosmeticC := &models.Category{
		Id:   uuid.NewV1(),
		Name: "Cosmetic",
	}
	server.InventoryService.CreateCategory(context.Background(), cosmeticC)

	// Taxes
	bst := &models.Tax{
		Id:        uuid.NewV1(),
		Name:      "Basic Sale Tax",
		Rate:      decimal.NewFromFloat32(10),
		Origin:    models.TaxOriginAll,
		Condition: models.ExemptToTax,
		Categories: map[uuid.UUID]bool{
			foodC.Id:    true,
			booksC.Id:   true,
			medicalC.Id: true,
		},
	}
	server.TaxService.CreateTax(context.Background(), bst)

	id := &models.Tax{
		Id:     uuid.NewV1(),
		Name:   "Import Duty",
		Rate:   decimal.NewFromFloat32(5),
		Origin: models.TaxOriginImport,
	}
	server.TaxService.CreateTax(context.Background(), id)

	// Case 1 items
	c1Book := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "book",
		CategoryId: booksC.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(12.49),
	}
	server.InventoryService.CreateItem(context.Background(), c1Book)

	c1CD := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "music CD",
		CategoryId: musicC.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(14.99),
	}
	server.InventoryService.CreateItem(context.Background(), c1CD)

	c1Chocolate := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "chocolate bar",
		CategoryId: foodC.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(0.85),
	}
	server.InventoryService.CreateItem(context.Background(), c1Chocolate)

	bid1, _ := server.SaleService.CreateBasket(context.Background())

	server.SaleService.AddItem(context.Background(), bid1, c1Book.Id, 2)
	server.SaleService.AddItem(context.Background(), bid1, c1CD.Id, 1)
	server.SaleService.AddItem(context.Background(), bid1, c1Chocolate.Id, 1)

	server.SaleService.CloseBasket(context.Background(), bid1)

	// Case 2 items
	c2Chocolate := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "box of chocolates ",
		CategoryId: foodC.Id,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(10.00),
	}
	server.InventoryService.CreateItem(context.Background(), c2Chocolate)

	c2Perfume := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "bottle of perfume",
		CategoryId: cosmeticC.Id,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(47.50),
	}
	server.InventoryService.CreateItem(context.Background(), c2Perfume)

	bid2, _ := server.SaleService.CreateBasket(context.Background())

	server.SaleService.AddItem(context.Background(), bid2, c2Chocolate.Id, 1)
	server.SaleService.AddItem(context.Background(), bid2, c2Perfume.Id, 1)

	server.SaleService.CloseBasket(context.Background(), bid2)

	// Case 3 items
	c3PerfumeImport := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "bottle of perfume",
		CategoryId: cosmeticC.Id,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(27.99),
	}
	server.InventoryService.CreateItem(context.Background(), c3PerfumeImport)

	c3PerfumeLocal := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "bottle of perfume",
		CategoryId: cosmeticC.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(18.99),
	}
	server.InventoryService.CreateItem(context.Background(), c3PerfumeLocal)

	c3Pills := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "packet of headache pills",
		CategoryId: medicalC.Id,
		Origin:     models.ItemOriginLocal,
		Price:      decimal.NewFromFloat32(9.75),
	}
	server.InventoryService.CreateItem(context.Background(), c3Pills)

	c3Chocolate := &models.InventoryItem{
		Id:         uuid.NewV1(),
		Name:       "imported chocolates",
		CategoryId: foodC.Id,
		Origin:     models.ItemOriginImported,
		Price:      decimal.NewFromFloat32(11.25),
	}
	server.InventoryService.CreateItem(context.Background(), c3Chocolate)

	bid3, _ := server.SaleService.CreateBasket(context.Background())

	server.SaleService.AddItem(context.Background(), bid3, c3PerfumeImport.Id, 1)
	server.SaleService.AddItem(context.Background(), bid3, c3PerfumeLocal.Id, 1)
	server.SaleService.AddItem(context.Background(), bid3, c3Chocolate.Id, 3)
	server.SaleService.AddItem(context.Background(), bid3, c3Pills.Id, 1)

	server.SaleService.CloseBasket(context.Background(), bid3)
}
