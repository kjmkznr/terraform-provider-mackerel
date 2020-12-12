package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateUrlPathWordDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diagnostics := stringDiag()(v, path)
		value, _ := v.(string)
		pattern := `^([A-Za-z0-9_][-A-Za-z0-9_]*)(/[A-Za-z0-9_][-A-Za-z0-9_]*)*$`
		if !regexp.MustCompile(pattern).MatchString(value) {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Value must be alphabet, number, '_', '-' or '/'",
				Detail:        fmt.Sprintf("value doesn't comply with restrictions (%q): %q", pattern, value),
				AttributePath: path,
			})
		}
		return diagnostics
	}
}

func validateMethodWordDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diagnostics := stringDiag()(v, path)
		value, _ := v.(string)
		if !stringInSlice(value, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}, false) {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Value must be HTTP Method",
				Detail:        fmt.Sprintf("value doesn't comply with restrictions (GET, POST, PUT, DELETE): %q", value),
				AttributePath: path,
			})
		}
		return diagnostics
	}
}

func validateChannelEvent(events []string, validEvents []string) error {
	for _, e := range events {
		if !stringInSlice(e, validEvents, true) {
			return fmt.Errorf("%s is not valid event", e)
		}
	}
	return nil
}

func stringDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diagnostics diag.Diagnostics
		if val, ok := v.(string); !ok {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Value must be string",
				Detail:        fmt.Sprintf("Value is not a string (type = %T)", val),
				AttributePath: path,
			})
		}
		return diagnostics
	}
}

func stringInSliceDiag(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diagnostics := stringDiag()(v, path)
		value, _ := v.(string)
		if len(diagnostics) == 0 {
			if !stringInSlice(value, valid, ignoreCase) {
				diagnostics = append(diagnostics, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Value has incorrect value",
					Detail:        fmt.Sprintf("expected value to be one of %v, got %s", strings.Join(valid, ","), value),
					AttributePath: path,
				})
			}
		}
		return diagnostics
	}
}

func stringInSlice(target string, valid []string, ignoreCase bool) bool {
	for _, s := range valid {
		if target == s || (ignoreCase && strings.ToLower(target) == strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func validateDiagFunc(validateFunc func(interface{}, string) ([]string, []error)) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		warnings, errs := validateFunc(i, fmt.Sprintf("%+v", path))
		var diagnostics diag.Diagnostics
		for _, warning := range warnings {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  warning,
			})
		}
		for _, err := range errs {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diagnostics
	}
}
