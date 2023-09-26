package helper

import "strings"

func SplitArgs(arg string) []string {
	b := []string{}
	for _, a := range strings.Split(arg, ",") {
		a = strings.Trim(a, " ")
		if len(a) > 0 {
			b = append(b, a)
		}
	}
	return b
}
