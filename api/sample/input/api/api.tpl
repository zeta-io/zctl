package api

import (
    "github.com/gin-gonic/gin"
    {{ range $package := .packages }}
        "{{- $package }}"
    {{ end }}
)

{{ range $index, $api := .schema.Apis }}
	func {{ apiFunc $api }}Api(c *gin.Context, req model.{{ apiFunc $api }}Req) (model.{{ apiResponse $api }}, error){
	    panic("not implements")
	}
{{ end }}