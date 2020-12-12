package provider

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func Test_validateMethodWord(t *testing.T) {
	tests := []struct {
		name           string
		arg            string
		wantWarnCount  int
		wantErrorCount int
	}{
		{
			name:           "Method GET",
			arg:            "GET",
			wantWarnCount:  0,
			wantErrorCount: 0,
		},
		{
			name:           "Method POST",
			arg:            "POST",
			wantWarnCount:  0,
			wantErrorCount: 0,
		},
		{
			name:           "Method DELETE",
			arg:            "DELETE",
			wantWarnCount:  0,
			wantErrorCount: 0,
		},
		{
			name:           "Method PUT",
			arg:            "PUT",
			wantWarnCount:  0,
			wantErrorCount: 0,
		},
		{
			name:           "Method HEAD",
			arg:            "HEAD",
			wantWarnCount:  0,
			wantErrorCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics := validateMethodWordDiag()(tt.arg, (cty.Path)(nil).Index(cty.StringVal("method")))
			gotWarns := 0
			gotErrors := 0
			for _, d := range diagnostics {
				if d.Severity == diag.Error {
					gotErrors++
				}
				if d.Severity == diag.Warning {
					gotWarns++
				}
			}
			if gotWarns != tt.wantWarnCount {
				t.Errorf("validateMethodWord() gotWarns = %v, want %v", gotWarns, tt.wantWarnCount)
			}
			if gotErrors != tt.wantErrorCount {
				t.Errorf("validateMethodWord() gotErrors = %v, want %v", gotErrors, tt.wantErrorCount)
			}
		})
	}
}
