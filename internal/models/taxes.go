package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
)

// TaxOrigin is defines tax scope according to inventory item origin
type TaxOrigin string

const (
	TaxOriginAll    TaxOrigin = "ALL"
	TaxOriginLocal  TaxOrigin = "LOCAL"
	TaxOriginImport TaxOrigin = "IMPORT"
)

// TaxCondition is defines category given in tax subject or exempt from tax.
type TaxCondition string

const (
	UnknownTC    TaxCondition = "UNKNOWN"
	ExemptToTax  TaxCondition = "EXEMPT"  // refers to only tax types in context will be free from tax
	SubjectToTax TaxCondition = "SUBJECT" // refers to only tax types in context will be effected from tax
)

// Tax
type Tax struct {
	Id     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Rate   decimal.Decimal `json:"rate"`
	Origin TaxOrigin       `json:"origin"`

	TaxScope
}

type TaxScope struct {
	Condition  TaxCondition `json:"condition"`
	Categories []uuid.UUID  `json:"categories"`
}

func (tt *TaxOrigin) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "LOCAL", "IMPORT", "ALL":
		*tt = TaxOrigin(str)

	default:
		*tt = TaxOriginAll
	}

	return nil
}

func (tc *TaxCondition) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "EXEMPT", "SUBJECT":
		*tc = TaxCondition(str)

	default:
		*tc = UnknownTC
	}

	return nil
}
