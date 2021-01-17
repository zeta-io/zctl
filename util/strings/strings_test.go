/**
2 * @Author: Nico
3 * @Date: 2021/1/18 2:05
4 */
package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapitalize(t *testing.T) {
	input := "abc"
	out := "Abc"
	assert.Equal(t, Capitalize(input), out)
}

func TestCamelCase(t *testing.T) {
	assert.Equal(t, CamelCase("/api/v1/users", '/'), "ApiV1Users")
	assert.Equal(t, CamelCase("api/v1/users", '/'), "apiV1Users")
}