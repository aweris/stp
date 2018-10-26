package inventory

import "errors"

var (
	ErrInvalidParameter = errors.New("invalid parameter")

	ErrInvalidCategoryId   = errors.New("invalid category id")
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryNotEmpty    = errors.New("category is not empty")

	ErrInvalidItemId    = errors.New("invalid item id")
	ErrInvalidItemName  = errors.New("invalid item name")
	ErrInvalidItemPrice = errors.New("invalid item price")
)
