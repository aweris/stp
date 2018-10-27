package sales

import "errors"

var (
	ErrInvalidParameter = errors.New("invalid parameter")

	ErrInvalidBasketId  = errors.New("invalid parameter")
	ErrInvalidItemCount = errors.New("invalid item count")
)
