package evaluator

import (
	"fmt"
	"testing"

	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/object"
	hparser "github.com/seblkma/go-himeji/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := hparser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		//{"42 >= 0", true}, // InfixExpression does not handle <= or >= yet
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
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

// GOFLAGS="-count=1" go test -run TestEvalInfixIntegerExpression
func TestEvalInfixIntegerExpression(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"3", 3},
		{"42", 42},
		{"-12", -12},
		{"-42", -42},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}

// GOFLAGS="-count=1" go test -run TestIfElseExpression
func TestIfElseExpression(t *testing.T) {
	// The ! operator negates the operand
	testInputs := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		expectedInt, ok := ti.expected.(int)
		fmt.Println(ti)
		if ok {
			if !testIntegerObject(t, evaluated, int64(expectedInt)) {
				break
			}
		} else {
			if !testNullObject(t, evaluated) {
				break
			}
		}
	}
}

// GOFLAGS="-count=1" go test -run TestReturnStatements
func TestReturnStatements(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"return 42;", 42},
		{"return 42; 9;", 42},
		{"return 2 * 21; 9;", 42},
		{"865; return 21 * 2; 911;", 42},
		// a nested if, expected to return 10
		{
			`
			if (10 > 1) {
			  if (10 > 1) {
			    return 10;
			  }
			  129
			  return 1;
			}
			`,
			10,
		},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}

// GOFLAGS="-count=1" go test -run TestErrorHandling
func TestErrorHandling(t *testing.T) {
	testInputs := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
		if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}
			return 1;
		}
		`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"uhoh;",
			"identifier not found: uhoh",
		},
	}

	for i, ti := range testInputs {
		evaluated := testEval(ti.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v) at test [%d]", evaluated, evaluated, i)
			continue
		}
		if errObj.Message != ti.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q at test [%d]", ti.expectedMessage, errObj.Message, i)
		}
	}
}

// GOFLAGS="-count=1" go test -run TestLetStatements
func TestLetStatements(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"let a = 42; a;", 42},
		{"let a = 5 * 5; a;", 25},
		{"let a = 42; let b = a; b;", 42},
		{"let a = 14; let b = a; let c = a + b + 14; c;", 42},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}
