package interperter

import (
	"bytes"
	"fmt"
	"github.com/zeta-io/zctl/api/schema"
	"github.com/zeta-io/zctl/api/types"
	strs "github.com/zeta-io/zctl/util/strings"
	"strings"
)

type ApiInterpreter interface {
	ApiFunc(api *schema.Api) string
	ApiQueries(api *schema.Api) []Entry
	ApiPathVariables(api *schema.Api) []Entry
	ApiBody(api *schema.Api) string
	ApiResponse(api *schema.Api) string
}

type ModelInterpreter interface {
	ModelName(model *schema.Model) string
	ModelFields(model *schema.Model) []int64
}

type Interpreter interface {
	ApiInterpreter
	ModelInterpreter
}

type Entry struct {
	Key string
	Value string
}

type ZetaInterpreter struct {}

func NewZeta() Interpreter{
	return &ZetaInterpreter{}
}

func (i *ZetaInterpreter) ApiFunc(api *schema.Api) string{
	funcName := strs.CamelCase(api.Path, '/')
	funcName = strs.Capitalize(api.Method) + funcName
	return strings.ReplaceAll(strings.ReplaceAll(funcName, "{", "_"), "}", "")
}

func (i *ZetaInterpreter) ApiQueries(api *schema.Api) []Entry{
	queries := make([]Entry, 0)
	for _, query := range api.Queries {
		queries = append(queries, Entry{
			Key: query.Name,
			Value: i.interpretTypes(query.Type),
		})
	}
	return queries
}

func (i *ZetaInterpreter) ApiPathVariables(api *schema.Api) []Entry{
	pathVariables := make([]Entry, 0)
	for _, pathVariable := range api.PathVariables {
		pathVariables = append(pathVariables, Entry{
			Key: pathVariable.Name,
			Value: i.interpretTypes(pathVariable.Type),
		})
	}
	return pathVariables
}

func (i *ZetaInterpreter) ApiBody(api *schema.Api) string{
	return i.interpretTypes(api.Body)
}

func (i *ZetaInterpreter) ApiResponse(api *schema.Api) string{
	return i.interpretTypes(api.Response)
}

func (i *ZetaInterpreter) ModelName(model *schema.Model) string{
	return model.Name
}

func (i *ZetaInterpreter) ModelFields(model *schema.Model) []int64{
	//fields := make([]Entry, 0)
	//for _, field := range model.Fields{
	//	fields = append(fields, Entry{
	//		Key: field.Name,
	//		Value: i.interpretTypes(field.Type),
	//	})
	//}
	//return fields
	return []int64{1, 2, 3}
}

func (i *ZetaInterpreter) interpretTypes(t *schema.Type) string{
	if t.Type.IsPrimitive(){
		return string(t.Type)
	}
	buff := bytes.Buffer{}
	if ! t.Required {
		buff.WriteRune('*')
	}
	if t.Type == types.Array{
		buff.WriteString("[]")
	}else if t.Type == types.Struct{
		buff.WriteString(t.Name)
	}else if t.Type == types.Map{
		buff.WriteString("[")
		buff.WriteString(i.interpretTypes(t.Key))
		buff.WriteString("]")
		buff.WriteString(i.interpretTypes(t.Value))
	}else if t.Type == types.Any{
		buff.WriteString("interface{}")
	}
	panic(fmt.Sprintf("type not support %v", t))
}