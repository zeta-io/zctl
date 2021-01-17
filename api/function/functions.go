package function

import (
	"github.com/zeta-io/zctl/api/types"
	"github.com/zeta-io/zctl/util/strings"
)

func Capitalize(str string) string {
	return strings.Capitalize(str)
}

func GoType(t types.Type) string {
	return types.GoTypes[t]
}
