package taxes

import "errors"

var (
	ErrInvalidParameter = errors.New("invalid parameter")

	ErrInvalidTaxId   = errors.New("invalid tax id")
	ErrInvalidTaxName = errors.New("invalid tax name")
	ErrInvalidTaxRate = errors.New("invalid tax rate")
)
