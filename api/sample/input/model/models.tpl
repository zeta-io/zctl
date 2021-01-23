package model

{{ range $index, $model := .schema.Models }}
	type {{modelName $model}} struct {
    {{ range $field := modelFields $model }}
        {{- capitalize $field.Key }} {{ $field.Value }}
    {{ end }}
	}
{{ end }}

{{ range $index, $api := .schema.Apis }}
	type {{ apiFunc $api }}Req struct {
        {{ range $queries := apiQueries $api }}
            {{- capitalize $queries.Key }} {{ $queries.Value }} `json:"{{ $queries.Key }}" param:"query,{{ $queries.Key }}"`
        {{ end }}
        {{ range $pathVariables := apiPathVariables $api }}
            {{- capitalize $pathVariables.Key }} {{ $pathVariables.Value }} `json:"{{ $pathVariables.Key }}" param:"path,{{ $pathVariables.Key }}"`
        {{ end }}
        {{ if $api.Body}}
            Body {{ apiBody $api }} `param:"body"`
        {{ end }}
	}
{{ end }}