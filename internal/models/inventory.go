package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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
	Imported   bool            `json:"imported"`
	Price      decimal.Decimal `json:"price"`
}
