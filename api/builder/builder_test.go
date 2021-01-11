package builder

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	schema2 "github.com/zeta-io/zctl/api/schema"
	"testing"
	"text/template"
)

const temp1 = `api list
{{ range $key, $value := .Apis }}
	{{ $value.Path}}
{{ end }}`

func TestTemplate(t *testing.T){
	schema, err := schema2.Parse(schema2.api1)
	assert.Equal(t, err, nil)
	bs, err := json.Marshal(schema)
	assert.Equal(t, err, nil)
	t.Log(string(bs))

	temp, err := template.New("test").Parse(temp1)
	assert.Equal(t, err, nil)
	buffer := bytes.Buffer{}
	err = temp.Execute(&buffer, schema)
	assert.Equal(t, err, nil)
	t.Log(buffer.String())
	t.Log("end")
}