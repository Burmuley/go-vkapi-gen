package main

import (
	"fmt"
	"strings"
)

type schemaError struct {
	errInfo string // additional text information of an error
	err     error  // original error
}

func (s schemaError) Error() string {
	var tmp string
	tmp = strings.ReplaceAll(s.errInfo, "\n", "")
	return fmt.Sprint(strings.Join([]string{"while parsing", tmp, "the following error occurred:", fmt.Sprintf("%s", s.err)}, " "))
}
