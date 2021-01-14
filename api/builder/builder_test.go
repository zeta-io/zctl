package builder

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-io/zctl/api/function"
	"github.com/zeta-io/zctl/api/schema"
	"io/ioutil"
	"testing"
	"text/template"
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

const temp1 = `api list
{{ range $key, $value := .Apis }}
	{{ $value.Path}}
{{ end }}`

func TestTemplate(t *testing.T){
	s, err := schema.Parse(api1)
	assert.Equal(t, err, nil)
	bs, err := json.Marshal(s)
	assert.Equal(t, err, nil)
	t.Log(string(bs))

	temp, err := template.New("test").Parse(temp1)
	assert.Equal(t, err, nil)
	buffer := bytes.Buffer{}
	err = temp.Execute(&buffer, s)
	assert.Equal(t, err, nil)
	t.Log(buffer.String())
	t.Log("end")
}

func TestTemplatesModelModels(t *testing.T){
	s, err := schema.Parse(api1)
	assert.Equal(t, err, nil)
	_, err = json.Marshal(s)
	assert.Equal(t, err, nil)

	b, err := ioutil.ReadFile("../sample/templates/model/models.tpl")
	assert.Equal(t, err, nil)

	temp, err := template.New("").Funcs(template.FuncMap{
		"capitalize": function.Capitalize,
	}).Parse(string(b))
	assert.Equal(t, err, nil)

	buffer := bytes.Buffer{}
	err = temp.Execute(&buffer, s)
	t.Log(err)
	assert.Equal(t, err, nil)
	t.Log(buffer.String())
	t.Log("end")
}