package models

import "github.com/satori/go.uuid"

type BasketState string

const (
	BasketStateOpened    BasketState = "OPENED"
	BasketStateClosed    BasketState = "CLOSED"
	BasketStateCancelled BasketState = "CANCELLED"
)

// Basket represents a record of the items that customer have chosen to buy
type Basket struct {
	Id    uuid.UUID        `json:"id"`
	Items map[SaleItem]int `json:"items"`
	State BasketState      `json:"state"`
}
