package helper

import "strings"

const connect = "#"

func IdFormat(s ...string) string {
	return strings.Join(s, connect)
}

func IdParse(s string) []string {
	return strings.Split(s, connect)
}
