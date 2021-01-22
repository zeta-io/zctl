package builder

import (
	"bytes"
	"github.com/zeta-io/zctl/api/function"
	"github.com/zeta-io/zctl/api/imports"
	"github.com/zeta-io/zctl/api/interperter"
	"github.com/zeta-io/zctl/api/schema"
	"github.com/zeta-io/zctl/errors"
	"github.com/zeta-io/zctl/util/file"
	"strings"
	"text/template"
)

type Builder struct {
	input  string
	output string
	schema *schema.Schema
	i interperter.Interpreter
}

func New(s *schema.Schema, i interperter.Interpreter, input, output string) (*Builder, error) {
	exist, err := file.IsExist(input)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.ErrBuildInputNotExist
	}
	return &Builder{
		schema: s,
		input:  input,
		output: output,
		i: i,
	}, nil
}

func (b *Builder) Generate() error {
	tpls, err := file.GetTpls(b.input)
	if err != nil {
		return err
	}
	outs := map[string]string{}

	packages := make([]string, 0)
	for _, tpl := range tpls{
		out := strings.Replace(tpl, b.input, b.output, 1)
		out = strings.Replace(out, ".tpl", ".go", 1)
		outs[tpl] = out

		outDir, err := file.GetDir(out)
		if err != nil{
			return err
		}
		if p, ok := imports.GoModuleRoot(outDir); ok{
			packages = append(packages, p)
		}
	}

	for _, tpl := range tpls {
		source, err := file.Read(tpl)
		if err != nil {
			return err
		}
		result, err := b.render(source, map[string]interface{}{
			"packages": packages,
			"schema": b.schema,
		})
		if err != nil {
			return err
		}
		bs, err := imports.Imports("", result)
		if err != nil{
			return err
		}
		err = file.Write(outs[tpl], bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) render(source []byte, data interface{}) ([]byte, error) {
	temp, err := template.New("").Funcs(template.FuncMap{
		"capitalize": function.Capitalize,
		"apiBody": b.i.ApiBody,
		"apiFunc": b.i.ApiFunc,
		"apiPathVariables": b.i.ApiPathVariables,
		"apiQueries": b.i.ApiQueries,
		"apiResponse": b.i.ApiResponse,
		"modelName": b.i.ModelName,
		"modelFields": b.i.ModelFields,
	}).Parse(string(source))
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	err = temp.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
