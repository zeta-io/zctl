/**
2 * @Author: Nico
3 * @Date: 2021/1/18 2:00
4 */
package strings

import (
	"bytes"
	"strings"
)

func Capitalize(str string) string {
	runes := []rune(str)
	if len(runes) > 0 {
		runes[0] = []rune(strings.ToUpper(string(runes[0:1])))[0]
	}
	return string(runes)
}

func CamelCase(str string, split rune) string{
	runes := []rune(str)
	n := false
	buf := bytes.Buffer{}
	for _, r := range runes{
		if r == split{
			n = true
			continue
		}
		if n{
			n = false
			buf.WriteString(strings.ToUpper(string(r)))
		}else{
			buf.WriteRune(r)
		}
	}
	return buf.String()
}