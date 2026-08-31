package evaluator

import (
	"testing"
)

func TestAssignment(t *testing.T) {
	input := `
var x = 5
x = 10
x
`
	testIntegerObject(t, testEval(input), 10)
}
