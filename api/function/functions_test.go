package function

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCapitalize(t *testing.T) {
	input := "abc"
	out := "Abc"
	assert.Equal(t, Capitalize(input), out)
}