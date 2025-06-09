package internal

import (
	"testing"
)

func TestSanitizeAccountName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Production Account", "production-account"},
		{"Test_Environment", "test-environment"},
		{"Dev Account 123", "dev-account-123"},
		{"Account-Name", "account-name"},
		{"  Spaced  Name  ", "spaced-name"},
		{"Special@#$Characters", "specialcharacters"},
		{"Mixed_Case-Name 123", "mixed-case-name-123"},
		{"", ""},
		{"---test---", "test"},
	}

	for _, test := range tests {
		result := sanitizeAccountName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeAccountName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}