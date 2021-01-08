package schema

import (
	"encoding/json"
	"testing"
)

const api1 = `api
  api/v1/users
    get ?page=uint&size=uint [userOutput]
    get /uid=uint userOutput
    get ?postUsersInput userOutput

model
  postUsersInput{age uint16, name string}
  userOutput
    id uint64
    age uint16
    name string
`

func TestParse(t *testing.T) {
	schema := Parse(api1)
	bs, err := json.Marshal(schema)
	t.Log(err)
	t.Log(string(bs))
}