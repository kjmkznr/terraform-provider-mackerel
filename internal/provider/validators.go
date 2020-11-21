package provider

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

func validateChannelEvent(events []string, validEvents []string) error {
	for _, e := range events {
		if !isStringInSlice(e, validEvents) {
			return fmt.Errorf("%s is not valid event", e)
		}
	}
	return nil
}

func isStringInSlice(target string, strings []string) bool {
	for _, s := range strings {
		if target == s {
			return true
		}
	}
	return false
}
