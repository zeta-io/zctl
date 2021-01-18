package model

{{ range $index, $model := .schema.Models }}
	type {{capitalize $model.Name}} struct {
    {{ range $index, $field := $model.Fields }}
        {{- capitalize $field.Name }} {{ goType $field.Type}}
    {{ end }}
	}
{{ end }}

{{ range $key, $value := .schema.Apis }}
	type {{ $value.Func }}Req struct {
        {{ range $field := $value.Queries }}
            {{- capitalize $field.Name }} {{ goType $field.Type}} `json:"{{ $field.Name }}" param:"query,{{ $field.Name }}"`
        {{ end }}
        {{ range $field := $value.PathVariables }}
            {{- capitalize $field.Name }} {{ goType $field.Type}} `json:"{{ $field.Name }}" param:"query,{{ $field.Name }}"`
        {{ end }}
        {{ if ne $value.Body ""}}
            Body {{ capitalize $value.Body }} `param:"body"`
        {{ end }}
	}
{{ end }}