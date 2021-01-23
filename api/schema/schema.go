package schema

import (
	"bytes"
	"fmt"
	"github.com/zeta-io/zctl/api/types"
	"github.com/zeta-io/zctl/errors"
	"github.com/zeta-io/zctl/util/stack"
	"strings"
)



type Schema struct {
	Apis   []*Api
	Models []*Model
}

type Model struct {
	Name   string
	Fields []*Field
}

type Api struct {
	Path          string
	Method        string
	Queries       []*Field
	PathVariables []*Field
	Body          *Type
	Response      *Type
}

type Field struct {
	Name string
	Type *Type
}

type Type struct {
	Name string
	Key *Type
	Value *Type
	Type types.Type
	Required bool
}

type reader struct {
	source []rune
	index  int
	buffer bytes.Buffer
	deep   int
}

func (a *Api) ID() string {
	return a.Method + ":" + a.Path
}

func Parse(input string) (*Schema, error) {
	r := &reader{source: []rune(input)}
	apis := make([]*Api, 0)
	models := make([]*Model, 0)

	lastDeep := 0
	status := 0
	stack := stack.New()
	for r.hasNext() {
		token, deep := r.next()
		if deep == 0 {
			if token == "api" {
				status = 1
			} else if token == "model" {
				status = 2
			}
			stack.Init()
			lastDeep = deep
			continue
		}
		dif := (deep - lastDeep) / 2
		if dif > 1 {
			return nil, errors.ErrSchemaFormatError
		}

		if status == 1 {
			if dif <= 0 {
				for i := 0; i <= -dif; i++ {
					stack.Pop()
				}
			}
			api, err := parseApi(token)
			if err != nil {
				return nil, err
			}
			if deep/2 == 1 {
				if api.Method != ""{
					apis = append(apis, api)
				}
				stack.Push(api)
			} else {
				cur := stack.Peek().(*Api)
				if api.Method != ""{
					api.Path = cur.Path + api.Path
					apis = append(apis, api)
				}
				stack.Push(api)
			}
		} else if status == 2 {
			if deep/2 == 1 {
				if dif <= 0 {
					for i := 0; i <= -dif; i++ {
						stack.Pop()
					}
				}
				model, err := parseModel(token)
				if err != nil {
					return nil, err
				}
				models = append(models, model)
				stack.Push(model)
			} else if deep/2 > 1 {
				cur := stack.Peek().(*Model)
				field, err := parseField(token)
				if err != nil {
					return nil, err
				}
				cur.Fields = append(cur.Fields, field)
			}
		}
		lastDeep = deep
	}
	return &Schema{
		Apis:   apis,
		Models: models,
	}, nil
}

func parseApi(token string) (*Api, error) {
	arr := strings.Split(token, " ")
	l := len(arr)

	method := ""
	path := ""
	var response *Type
	switch l {
	case 1:
		path = arr[0]
	case 2:
		method = arr[0]
		path = arr[1]
	case 3:
		method = arr[0]
		path = arr[1]

		t, err := ParseType(arr[2])
		if err != nil {
			return nil, err
		}
		response = t
	default:
		return nil, fmt.Errorf("not support api definition: %s", token)
	}

	var queries []*Field
	var body *Type
	if i := strings.Index(path, "?"); i > -1 {
		parameters := strings.Split(path[i+1:], "&")
		for _, parameter := range parameters {
			elements := strings.Split(parameter, "=")
			if len(elements) == 1 && body == nil {
				t, err := ParseType(elements[0])
				if err != nil {
					return nil, err
				}
				body = t
			} else if len(elements) == 2 {
				t, err := ParseType(elements[1])
				if err != nil {
					return nil, err
				}
				queries = append(queries, &Field{
					Name: elements[0],
					Type: t,
				})
			}
		}
		path = path[:i]
	}

	var pathVariables []*Field
	pathSegments := strings.Split(path, "/")
	path = ""
	for _, segment := range pathSegments {
		if segment == "" {
			continue
		}
		if i := strings.Index(segment, "="); i > -1 {
			elements := strings.Split(segment, "=")
			t, err := ParseType(elements[1])
			if err != nil {
				return nil, err
			}
			pathVariables = append(pathVariables, &Field{
				Name: elements[0],
				Type: t,
			})
			path += "/" + fmt.Sprintf("{%d}", len(pathVariables)-1)
		} else {
			path += "/" + segment
		}
	}
	if path == "" {
		path = "/"
	}
	return &Api{
		Path:          path,
		Method:        method,
		Queries:       queries,
		PathVariables: pathVariables,
		Body:          body,
		Response:      response,
	}, nil
}

func parseModel(token string) (*Model, error) {
	if strings.Contains(token, "{") {
		arr := strings.Split(token, "{")
		fields, err := parseFieldInline(strings.ReplaceAll(arr[1], "}", ""))
		return &Model{
			Name:   arr[0],
			Fields: fields,
		}, err
	} else {
		return &Model{
			Name: token,
		}, nil
	}
}

func parseFieldInline(token string) ([]*Field, error) {
	arr := strings.Split(token, ",")
	fields := make([]*Field, 0)
	for _, elem := range arr {
		targets := strings.Split(strings.TrimSpace(elem), " ")
		t, err := ParseType(targets[1])
		if err != nil {
			return nil, err
		}
		fields = append(fields, &Field{
			Name: targets[0],
			Type: t,
		})
	}
	return fields, nil
}

// ParseType to parse type form token:
func ParseType(token string) (*Type, error) {
	t := &Type{
		Required: true,
	}
	if token[0] == '*'{
		t.Required = false
		token = token[1:]
	}
	buff := bytes.Buffer{}
	stat := 0
	sign := 0
	for _, s := range token{
		switch s {
		case '[':
			if stat == 0{
				if buff.String() == "map"{
					t.Type = types.Map
					stat = 1
				}else{
					t.Type = types.Array
					stat = 3
				}
				buff.Reset()
			}else if stat == 1{
				sign ++
			}else{
				buff.WriteRune(s)
			}
		case ']':
			if stat == 1{
				var err error
				t.Key, err = ParseType(buff.String())
				if err != nil{
					return nil, err
				}
				stat = 2
				buff.Reset()
			}else if stat == 3{
				if buff.Len() > 0{
					return nil, errors.ErrTypesFormat
				}
				stat = 2
			}else{
				buff.WriteRune(s)
			}
		default:
			buff.WriteRune(s)
		}
	}
	if t.Type == types.Map || t.Type == types.Array{
		var err error
		t.Value, err = ParseType(buff.String())
		if err != nil{
			return nil, err
		}
	}else if token == "any"{
		t.Type = types.Any
	}else if ty, err := types.ParsePrimitive(token); err == nil{
		t.Type = ty
	}else{
		t.Type = types.Struct
		t.Name = token
	}
	return t, nil
}

func parseField(token string) (*Field, error) {
	arr := strings.Split(token, " ")
	t, err := ParseType(arr[1])
	if err != nil {
		return nil, err
	}
	return &Field{
		Name: arr[0],
		Type: t,
	}, nil
}

func (r *reader) hasNext() bool {
	return r.index < len(r.source)-1
}

func (r *reader) next() (string, int) {
	ready := false
	deep := 0
	r.foreach(func(c rune) bool {
		if c == '\n' {
			return false
		}
		if !ready {
			if c == ' ' {
				deep += 1
			} else if c == '\t' {
				deep += 2
			} else {
				ready = true
			}
		}
		if ready {
			r.buffer.WriteRune(c)
		}
		return true
	})
	if deep%2 == 1 {
		panic("format err: 001")
	}
	token := r.buffer.String()
	r.buffer.Reset()
	return token, deep
}

func (r *reader) foreach(handle func(c rune) bool) {
	for r.index < len(r.source) {
		isBreak := !handle(r.source[r.index])
		r.index++
		if isBreak {
			break
		}
	}
}
