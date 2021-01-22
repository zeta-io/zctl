package api

import (
    "github.com/gin-gonic/gin"
    {{ range $package := .packages }}
        "{{- $package }}"
    {{ end }}
)

{{ range $key, $value := .schema.Apis }}
	func {{ $value.Func }}(c *gin.Context, req model.{{ $value.Func }}Req) (model.{{ capitalize $value.Response }}, error){
	    panic("not implements")
	}
{{ end }}