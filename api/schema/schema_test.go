package schema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const api1 = `api
  api/v1/users
    get ?page=uint&size=uint [userOutput]
    get /uid=uint userOutput
    post ?postUsersInput userOutput

model
  postUsersInput{age uint16, name string}
  userOutput
    id uint64
    age uint16
    name string
`

func TestParse(t *testing.T) {
	schema, err := Parse(api1)
	assert.Equal(t, err, nil)
	bs, err := json.Marshal(schema)
	assert.Equal(t, err, nil)
	t.Log(string(bs))
}

func TestParseType(t *testing.T) {
	token := "map[string]string"
	ty, err := ParseType(token)
	assert.Equal(t, err, nil)

	str, _ := json.Marshal(ty)
	t.Log(string(str))
}