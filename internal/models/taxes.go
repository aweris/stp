package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
)

type TaxCondition string
type TaxType string

const (
	UnknownTC    TaxCondition = "UNKNOWN"
	ExemptToTax  TaxCondition = "EXEMPT"  // refers to only tax types in context will be free from tax
	SubjectToTax TaxCondition = "SUBJECT" // refers to only tax types in context will be effected from tax

	UnknownTaxType TaxType = "UNKNOWN"
	CategoryTax    TaxType = "CATEGORY" // refers to items will be filtered by category
	OriginTax      TaxType = "ORIGIN"   // refers to items will be filtered by item origin
)

// Tax
type Tax struct {
	Id        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Rate      decimal.Decimal `json:"rate"`
	Type      TaxType         `json:"type"`
	Condition TaxCondition    `json:"condition"`
	Context   []uuid.UUID     `json:"context"`
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

func (tt *TaxType) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "CATEGORY", "ORIGIN":
		*tt = TaxType(str)

	default:
		*tt = UnknownTaxType
	}

	return nil
}
