package utils

import (
	"fmt"
)

//Block puts the given text in a single-line code block
func Block(s interface{}) string {
	return fmt.Sprintf("`%v`", s)
}

//Bold bolds the given text
func Bold(s interface{}) string {
	return fmt.Sprintf("**%v**", s)
}
