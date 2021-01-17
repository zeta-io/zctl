package api

{{ range $key, $value := .Apis }}
	func {{ $value.Func }}(){}
{{ end }}