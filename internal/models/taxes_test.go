package models_test

import (
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaxCondition_UnmarshalText(t *testing.T) {
	str := "{\"condition\":\"EXEMPT\"}"

	var tax models.Tax

	json.Unmarshal([]byte(str), &tax)

	assert.Equal(t, tax.Condition, models.ExemptToTax)
}

func TestTaxCondition_UnmarshalText_WhenNotKnownOption_ThenShouldAssignUnknown(t *testing.T) {
	str := "{\"condition\":\"NOT OPTION\"}"

	var tax models.Tax

	json.Unmarshal([]byte(str), &tax)

	assert.Equal(t, tax.Condition, models.UnknownTC)
}

func TestTaxOrigin_UnmarshalText(t *testing.T) {
	str := "{\"origin\":\"LOCAL\"}"

	var tax models.Tax

	json.Unmarshal([]byte(str), &tax)

	assert.Equal(t, tax.Origin, models.TaxOriginLocal)
}

func TestTaxType_UnmarshalText_WhenNotKnownOption_ThenShouldAssignUnknown(t *testing.T) {
	str := "{\"origin\":\"NOT OPTION\"}"

	var tax models.Tax

	json.Unmarshal([]byte(str), &tax)

	assert.Equal(t, tax.Origin, models.TaxOriginAll)
}
