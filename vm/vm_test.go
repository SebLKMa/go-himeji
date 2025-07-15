package vm

import (
	"fmt"
	"testing"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/compiler"
	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/object"
	"github.com/seblkma/go-himeji/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch exp := expected.(type) {
	case int:
		err := testIntegerObject(int64(exp), actual)
		if err != nil {
			t.Errorf("testExpectedObject failed: %s", err)
		}
	}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.ByteCode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

// GOFLAGS="-count=1" go test -run TestIntegerArithmetic
func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3}, // TODO: expected 3, but for now only has 2 at top of the stack
	}

	runVmTests(t, tests)
}
