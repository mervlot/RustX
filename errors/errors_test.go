package errors

import (
	"testing"
)

func TestErrorFormatting(t *testing.T) {
	synErr := NewSyntaxError("Expected ')' after expression", 5, 12)
	expectedSyn := "SyntaxError: Expected ')' after expression at line 5 column 12"
	if synErr.Error() != expectedSyn {
		t.Errorf("expected %q, got %q", expectedSyn, synErr.Error())
	}

	runErr := NewRuntimeError("Division by zero", 10, 4)
	expectedRun := "RuntimeError: Division by zero at line 10 column 4"
	if runErr.Error() != expectedRun {
		t.Errorf("expected %q, got %q", expectedRun, runErr.Error())
	}
}
