package inventory

import "errors"

var (
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrInvalidCategoryId   = errors.New("invalid category id")
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryNotEmpty    = errors.New("category is not empty")
)
