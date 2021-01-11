package api

{{ range $key, $value := .Apis }}
	{{ $value.Path}}
{{ end }}`