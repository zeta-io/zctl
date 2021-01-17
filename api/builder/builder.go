package builder

import (
	"bytes"
	"github.com/zeta-io/zctl/api/function"
	"github.com/zeta-io/zctl/api/schema"
	"github.com/zeta-io/zctl/errors"
	"github.com/zeta-io/zctl/util/file"
	"golang.org/x/tools/imports"
	"strings"
	"text/template"
)

type Builder struct {
	input  string
	output string
	schema *schema.Schema
}

func New(s *schema.Schema, input, output string) (*Builder, error) {
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
	}, nil
}

func (b *Builder) Generate() error {
	tpls, err := file.ReadTpls(b.input)
	if err != nil {
		return err
	}
	outs := make([]string, 0)
	for _, tpl := range tpls {
		source, err := file.Read(tpl)
		if err != nil {
			return err
		}
		result, err := b.render(source)
		if err != nil {
			return err
		}
		out := strings.Replace(tpl, b.input, b.output, 1)
		out = strings.Replace(out, ".tpl", ".go", 1)
		bs, err := imports.Process(out, result, &imports.Options{FormatOnly: false, Comments: true, TabIndent: true, TabWidth: 8})
		if err != nil{
			return err
		}
		err = file.Write(out, bs)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}
	return nil
}

func (b *Builder) render(source []byte) ([]byte, error) {
	temp, err := template.New("").Funcs(template.FuncMap{
		"capitalize": function.Capitalize,
		"goType":     function.GoType,
	}).Parse(string(source))
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	err = temp.Execute(&buf, b.schema)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
