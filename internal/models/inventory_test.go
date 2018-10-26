package models_test

import (
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemOrigin_UnmarshalText(t *testing.T) {
	str := "{\"origin\":\"IMPORT\"}"

	var item models.InventoryItem

	json.Unmarshal([]byte(str), &item)

	assert.Equal(t, item.Origin, models.ItemOriginImported)
}

func TestItemOrigin_UnmarshalText_WhenNotKnownOption_ThenShouldAssignUnknown(t *testing.T) {
	str := "{\"origin\":\"NOT OPTION\"}"

	var item models.InventoryItem

	json.Unmarshal([]byte(str), &item)

	assert.Equal(t, item.Origin, models.ItemOriginLocal)
}