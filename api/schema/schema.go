package schema

import (
	"bytes"
	"fmt"
	"github.com/zeta-io/zctl/api/types"
	"github.com/zeta-io/zctl/errors"
	"github.com/zeta-io/zctl/util/stack"
	strs "github.com/zeta-io/zctl/util/strings"
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
	Func		  string
	Method        string
	Queries       []*Field
	PathVariables []*Field
	Body          string
	Response      string
}

type Field struct {
	Name string
	Type types.Type
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
				apis = append(apis, api)
				stack.Push(api)
			} else {
				cur := stack.Peek().(*Api)
				api.Path = cur.Path + api.Path
				apis = append(apis, api)
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

	method := "any"
	path := ""
	response := ""
	switch l {
	case 1:
		path = arr[0]
	case 2:
		method = arr[0]
		path = arr[1]
	case 3:
		method = arr[0]
		path = arr[1]
		response = arr[2]
	default:
		return nil, fmt.Errorf("not support api definition: %s", token)
	}

	var queries []*Field
	var body string
	if i := strings.Index(path, "?"); i > -1 {
		parameters := strings.Split(path[i+1:], "&")
		for _, parameter := range parameters {
			elements := strings.Split(parameter, "=")
			if len(elements) == 1 && body == "" {
				body = elements[0]
			} else if len(elements) == 2 {
				t, err := types.Parse(elements[1])
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
			t, err := types.Parse(elements[1])
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
	funcName := strs.CamelCase(path, '/')
	funcName = strs.Capitalize(method) + funcName
	funcName = strings.ReplaceAll(strings.ReplaceAll(funcName, "{", "_"), "}", "")
	return &Api{
		Path:          path,
		Func: 		   funcName,
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
		t, err := types.Parse(targets[1])
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

func parseField(token string) (*Field, error) {
	arr := strings.Split(token, " ")
	t, err := types.Parse(arr[1])
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
