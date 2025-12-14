package naming

import "strings"

func ToKebab(s string, options ...Option) string {
	return strings.ReplaceAll(ToSnake(s, options...), "_", "-")
}
