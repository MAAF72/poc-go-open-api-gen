package main

import "regexp"

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// ClearString :nodoc
func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
