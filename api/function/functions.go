package function

import "strings"

func Capitalize(str string) string{
	runes := []rune(str)
	if len(runes) > 0{
		runes[0] = []rune(strings.ToUpper(string(runes[0:1])))[0]
	}
	return string(runes)
}