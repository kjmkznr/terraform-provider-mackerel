package provider

import (
	"reflect"
	"testing"
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
			gotWarns, gotErrors := validateMethodWord(tt.arg, "method")
			if !reflect.DeepEqual(len(gotWarns), tt.wantWarnCount) {
				t.Errorf("validateMethodWord() gotWarns = %v, want %v", gotWarns, tt.wantWarnCount)
			}
			if !reflect.DeepEqual(len(gotErrors), tt.wantErrorCount) {
				t.Errorf("validateMethodWord() gotErrors = %v, want %v", gotErrors, tt.wantErrorCount)
			}
		})
	}
}
