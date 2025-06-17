package evaluator

import (
	"testing"

	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/object"
	hparser "github.com/seblkma/go-himeji/parser"
)

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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

// GOFLAGS="-count=1" go test -run TestEvalIntegerExpression
func TestEvalIntegerExpression(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"3", 3},
		{"42", 42},
		{"-12", -12},
		{"-42", -42},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}

// GOFLAGS="-count=1" go test -run TestEvalBooleanExpression
func TestEvalBooleanExpression(t *testing.T) {
	testInputs := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testBooleanObject(t, evaluated, ti.expected)
	}
}

// GOFLAGS="-count=1" go test -run TestBangOperator
func TestBangOperator(t *testing.T) {
	// The ! operator negates the operand
	testInputs := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!42", false},
		{"!!true", true},   // negates true to false then negates false to true
		{"!!false", false}, // vice-versa above
		{"!!42", true},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testBooleanObject(t, evaluated, ti.expected)
	}
}
