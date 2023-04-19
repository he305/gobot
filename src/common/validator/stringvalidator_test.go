package validator_test

import (
	"gobot/src/common/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyShouldReturnFalse(t *testing.T) {
	str := ""

	expected := false

	actual := validator.ValidateString(str)

	assert.Equal(t, expected, actual)
}

func TestSpacesShouldReturnFalse(t *testing.T) {
	str := "\t\n"

	expected := false

	actual := validator.ValidateString(str)

	assert.Equal(t, expected, actual)
}

func TestOkShouldReturnTrue(t *testing.T) {
	str := "1\t\n"

	expected := true

	actual := validator.ValidateString(str)

	assert.Equal(t, expected, actual)
}
