package mackerel

import (
	"fmt"
	"regexp"
)

func validateUrlPathWord(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	pattern := `^([A-Za-z0-9_][-A-Za-z0-9_]*)(/[A-Za-z0-9_][-A-Za-z0-9_]*)*$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}
