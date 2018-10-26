package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type ItemOrigin string

const (
	ItemOriginImported ItemOrigin = "imported"
	ItemOriginLocal    ItemOrigin = "local"
)

// Category is used to define if an item is suitable for billing or pricing.
type Category struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// InventoryItem represents the physical good available
type InventoryItem struct {
	Id         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	CategoryId uuid.UUID       `json:"category"`
	Origin     ItemOrigin      `json:"origin"`
	Price      decimal.Decimal `json:"price"`
}
