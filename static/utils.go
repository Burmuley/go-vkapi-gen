package go_vkapi

import (
	"fmt"
	"strings"
)

// SliceToString converts any slice to a string with slice elements comma delimited
func SliceToString(slice interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slice)), ","), "[]")
}
