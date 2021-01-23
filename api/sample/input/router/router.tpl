package router

import (
    "github.com/gin-gonic/gin"
    "github.com/zeta-io/ginx"
    "github.com/zeta-io/zeta"
    {{ range $package := .packages }}
        "{{- $package }}"
    {{ end }}
)

func Start(addr ...string) error{
    router := zeta.Router("")
    {{- range $index, $api := .schema.Apis }}
        router.{{- capitalize $api.Method}}("{{ apiPath $api }}", api.{{apiFunc $api}}Api)
    {{- end }}
    return zeta.New(router, ginx.New(gin.New())).Run(addr...)
}
