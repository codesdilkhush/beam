package backend

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
		success  bool
	}{
		{"2 + 2", "4", true},
		{"10 * 5 - 2", "48", true},
		{"sqrt(144)", "12", true},
		{"abs(-42)", "42", true},
		{"10 / 4", "2.5", true},
		{"hello world", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		res := Evaluate(tt.expr)
		if res.Success != tt.success {
			t.Errorf("Evaluate(%q).Success = %v; want %v", tt.expr, res.Success, tt.success)
		}
		if tt.success && res.Result != tt.expected {
			t.Errorf("Evaluate(%q).Result = %q; want %q", tt.expr, res.Result, tt.expected)
		}
	}
}
