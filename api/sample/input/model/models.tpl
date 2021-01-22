package model

{{ range $index, $model := .schema.Models }}
	type {{modelName $model}} struct {
    {{ range $key, $value := modelFields $model }}
        {{- $key }} {{ $value}}
    {{ end }}
	}
{{ end }}

{{ range $index, $api := .schema.Apis }}
	type {{ apiFunc $api }}Req struct {
        {{ range $key, $value := apiQueries $api }}
            {{- capitalize $key }} {{ $value }} `json:"{{ $key }}" param:"query,{{ $key }}"`
        {{ end }}
        {{ range $key, $value := apiPathVariables $api }}
            {{- capitalize $key }} {{ $value }} `json:"{{ $key }}" param:"path,{{ $key }}"`
        {{ end }}
        {{ if ne $api.Body ""}}
            Body {{ capitalize apiBody $api }} `param:"body"`
        {{ end }}
	}
{{ end }}