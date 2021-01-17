package model

{{ range $index, $model := .Models }}
	type {{capitalize $model.Name}} struct {
    {{ range $index, $field := $model.Fields }}
        {{- capitalize $field.Name }} {{ goType $field.Type}}
    {{ end }}
	}
{{ end }}