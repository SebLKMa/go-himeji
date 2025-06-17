package evaluator

import (
	"testing"

	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/object"
	hparser "github.com/seblkma/go-himeji/parser"
)

// GOFLAGS="-count=1" go test -run TestEvalIntegerExpression
func TestEvalIntegerExpression(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"3", 3},
		{"42", 42},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := hparser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}
