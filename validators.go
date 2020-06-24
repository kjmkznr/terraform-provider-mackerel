package mackerel

import (
	"fmt"
	"net/http"
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

func validateMethodWord(v interface{}, k string) (warns []string, errors []error) {
	value := v.(string)
	switch value {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		return
	default:
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (GET, POST, PUT, DELETE): %q",
			k, value))
	}
	return
}
