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
	ApiPath(api *schema.Api) string
	ApiFunc(api *schema.Api) string
	ApiQueries(api *schema.Api) []Entry
	ApiPathVariables(api *schema.Api) []Entry
	ApiBody(api *schema.Api) string
	ApiResponse(api *schema.Api) string
}

type ModelInterpreter interface {
	ModelName(model *schema.Model) string
	ModelFields(model *schema.Model) []Entry
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

func (i *ZetaInterpreter) ApiPath(api *schema.Api) string{
	if len(api.PathVariables) == 0{
		return api.Path
	}
	path := api.Path
	for i, pathVariable := range api.PathVariables{
		path = strings.Replace(path, fmt.Sprintf("{%d}", i), pathVariable.Name, 1)
	}
	return path
}

func (i *ZetaInterpreter) ApiFunc(api *schema.Api) string{
	funcName := strs.CamelCase(api.Path, '/')
	funcName = strs.Capitalize(api.Method) + funcName
	return strings.ReplaceAll(strings.ReplaceAll(funcName, "{", "_"), "}", "_")
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
	return strs.Capitalize(model.Name)
}

func (i *ZetaInterpreter) ModelFields(model *schema.Model) []Entry{
	fields := make([]Entry, 0)
	for _, field := range model.Fields{
		fields = append(fields, Entry{
			Key: field.Name,
			Value: i.interpretTypes(field.Type),
		})
	}
	return fields
}

func (i *ZetaInterpreter) interpretTypes(t *schema.Type) string{
	buff := bytes.Buffer{}
	if ! t.Required {
		buff.WriteRune('*')
	}
	if t.Type.IsPrimitive(){
		switch t.Type {
		case types.Int: buff.WriteString("int")
		case types.Int8: buff.WriteString("int8")
		case types.Int16: buff.WriteString("int16")
		case types.Int32: buff.WriteString("int32")
		case types.Int64: buff.WriteString("int64")
		case types.UInt: buff.WriteString("uint")
		case types.UInt8: buff.WriteString("uint8")
		case types.UInt16: buff.WriteString("uint16")
		case types.UInt32: buff.WriteString("uint32")
		case types.UInt64: buff.WriteString("uint64")
		case types.Bool: buff.WriteString("bool")
		case types.Float32: buff.WriteString("float32")
		case types.Float64: buff.WriteString("float64")
		case types.Time: buff.WriteString("time.Time")
		case types.String: buff.WriteString("string")
		}
	}else if t.Type == types.Array{
		buff.WriteString("[]")
	}else if t.Type == types.Struct{
		buff.WriteString(strs.Capitalize(t.Name))
	}else if t.Type == types.Map{
		buff.WriteString("[")
		buff.WriteString(i.interpretTypes(t.Key))
		buff.WriteString("]")
		buff.WriteString(i.interpretTypes(t.Value))
	}else if t.Type == types.Any{
		buff.WriteString("interface{}")
	}else{
		panic(fmt.Sprintf("type not support %v", t))
	}
	return buff.String()
}