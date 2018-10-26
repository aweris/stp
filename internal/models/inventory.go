package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
)

type ItemOrigin string

const (
	ItemOriginUnknown  ItemOrigin = "UNKNOWN"
	ItemOriginImported ItemOrigin = "IMPORT"
	ItemOriginLocal    ItemOrigin = "LOCAL"
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

func (io *ItemOrigin) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "IMPORT", "LOCAL":
		*io = ItemOrigin(str)

	default:
		*io = ItemOriginUnknown
	}

	return nil
}
