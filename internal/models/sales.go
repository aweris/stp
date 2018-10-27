package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type BasketState string

const (
	BasketStateOpened    BasketState = "OPENED"
	BasketStateClosed    BasketState = "CLOSED"
	BasketStateCancelled BasketState = "CANCELLED"
)

// Basket represents a record of the items that customer have chosen to buy
type Basket struct {
	Id    uuid.UUID                 `json:"id"`
	Items map[uuid.UUID]*BasketItem `json:"items"`
	State BasketState               `json:"state"`
}

type BasketItem struct {
	*SaleItem

	Count int `json:"count"`
}

// Receipt represents written acknowledgment that something of value has been received.
type Receipt struct {
	Id uuid.UUID `json:"id"`

	basketID   uuid.UUID
	Items      []*BasketItem   `json:"items"`
	TotalTax   decimal.Decimal `json:"total_tax"`
	TotalPrice decimal.Decimal `json:"total_price"`
	TotalGross decimal.Decimal `json:"total_gross"`
}
