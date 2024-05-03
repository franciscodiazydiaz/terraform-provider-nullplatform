package nullplatform

import (
	"fmt"
	"testing"
)

func TestParseNrn(t *testing.T) {
	tests := []struct {
		nrn        string
		queryParam string
		expected   string
		wantErr    bool
	}{
		{"organization=1:account=2:namespace=3:application=4:scope=5", "application", "organization=1:account=2:namespace=3:application=4", false},
		{"organization=1:account=2:namespace=3", "organization", "organization=885129321", false},
		{"organization=1:account=2:namespace=3", "namespace", "organization=885129321:account=683868216:namespace=1764294975", false},
		{"organization=1:account=2:namespace=3", "account", "", true},
		{"", "organization", "", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("input=%s,queryParam=%s", test.nrn, test.queryParam), func(t *testing.T) {
			got, err := ParseNrn(test.nrn, test.queryParam)
			if (err != nil) != test.wantErr {
				t.Errorf("ParseNrn() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.expected {
				t.Errorf("ParseString() = %v, want %v", got, test.expected)
			}
		})
	}
}
