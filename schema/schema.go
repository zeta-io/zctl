package schema

import (
	"bytes"
	"fmt"
	"github.com/zeta-io/zctl/util/stack"
	"strings"
)

type Schema struct {
	Apis map[string]*Api
	Models map[string]*Model
}

type Model struct {
	Name string
	Fields []Field
}

type Api struct {
	Path string
	Method string
	Queries []Field
	Body string
	Return string
	Subs map[string]*Api
}

type Field struct {
	Name string
	Type string
}

type reader struct {
	source []rune
	index int
	buffer bytes.Buffer
	deep int
}

func Parse(input string) *Schema{
	r := &reader{source: []rune(input)}
	apis := map[string]*Api{}
	models := map[string]*Model{}

	lastDeep := 0
	status := 0
	stack := stack.New()
	for r.hasNext(){
		token, deep := r.next()
		if deep == 0{
			if token == "api"{
				status = 1
			}else if token == "model"{
				status = 2
			}
			stack.Init()
			lastDeep = deep
			continue
		}
		dif := (deep - lastDeep) / 2
		if dif > 1{
			panic("format err: 002")
		}

		if status == 1{
			if dif <= 0{
				for i := 0; i <= -dif; i ++{
					stack.Pop()
				}
			}
			api := parseApi(token)
			if deep/2 == 1{
				apis[api.Path] = api
				stack.Push(api)
			}else{
				cur := stack.Peek().(*Api)
				cur.Subs[api.Path] = api
				stack.Push(api)
			}
		}else if status == 2{
			if deep/2 == 1{
				if dif <= 0{
					for i := 0; i <= -dif; i ++{
						stack.Pop()
					}
				}

				model := parseModel(token)
				models[model.Name] = model
				stack.Push(model)
			}else if deep/2 > 1{
				cur := stack.Peek().(*Model)
				cur.Fields = append(cur.Fields, parseField(token))
			}
		}
		lastDeep = deep
	}
	return &Schema{
		Apis: apis,
		Models: models,
	}
}

func parseApi(token string) *Api{
	arr := strings.Split(token, " ")
	l := len(arr)
	if l == 1{
		return &Api{
			Path: arr[0],
			Subs: map[string]*Api{},
		}
	}else if l == 2{
		return &Api{
			Path: arr[1],
			Method: arr[0],
			Subs: map[string]*Api{},
		}
	}else if l == 3{
		return &Api{
			Path: arr[1],
			Method: arr[0],
			Return: arr[2],
			Subs: map[string]*Api{},
		}
	}
	panic(fmt.Sprintf("not support api parse: %s", token))
}

func parseModel(token string) *Model{
	if strings.Contains(token, "{") {
		arr := strings.Split(token, "{")
		return &Model{
			Name: arr[0],
			Fields: parseFieldInline(strings.ReplaceAll(arr[1], "}", "")),
		}
	}else{
		return &Model{
			Name: token,
		}
	}
}

func parseFieldInline(token string) []Field{
	arr := strings.Split(token, ",")
	fields := make([]Field, 0)
	for _, elem := range arr{
		targets := strings.Split(strings.TrimSpace(elem), " ")
		fields = append(fields, Field{
			Name: targets[0],
			Type: targets[1],
		})
	}
	return fields
}

func parseField(token string) Field{
	arr := strings.Split(token, " ")
	return Field{
		Name: arr[0],
		Type: arr[1],
	}
}

func (r *reader) hasNext() bool{
	return r.index < len(r.source) - 1
}

func (r *reader) next() (string, int){
	ready := false
	deep := 0
	r.foreach(func(c rune) bool {
		if c == '\n'{
			return false
		}
		if ! ready{
			if c == ' '{
				deep += 1
			}else if c == '\t'{
				deep += 2
			}else{
				ready = true
			}
		}
		if ready{
			r.buffer.WriteRune(c)
		}
		return true
	})
	if deep % 2 == 1{
		panic("format err: 001")
	}
	token := r.buffer.String()
	r.buffer.Reset()
	return token, deep
}

func (r *reader) foreach(handle func(c rune) bool){
	for ; r.index < len(r.source); {
		isBreak := ! handle(r.source[r.index])
		r.index ++
		if isBreak {
			break
		}
	}
}

