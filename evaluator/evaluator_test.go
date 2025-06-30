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
		{
			`"guten" - "tag!"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
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

// GOFLAGS="-count=1" go test -run TestFunctionObject
func TestFunctionObject(t *testing.T) {
	testBody := "{ x + 2; };"
	testInput := "fn(x) " + testBody
	expectedBody := "(x + 2)"

	evaluated := testEval(testInput)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

// GOFLAGS="-count=1" go test -run TestFunctionApplication
func TestFunctionApplication(t *testing.T) {
	testInputs := []struct {
		input    string
		expected int64
	}{
		{"let identityFn = fn(x) { x; }; identityFn(5);", 5},
		{"let identityFn = fn(x) { return x; }; identityFn(5);", 5},
		{"let doubleFn = fn(x) { x * 2; }; doubleFn(5);", 10},
		{"let addFn = fn(x, y) { x + y; }; addFn(20, 22);", 42},
		{"let addFn = fn(x, y) { x + y; }; addFn(10 + 10, addFn(20, 2));", 42}, // need to pass 20, 22 to outer addFn
		{"fn(x) { x; }(42)", 42},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		testIntegerObject(t, evaluated, ti.expected)
	}
}

// GOFLAGS="-count=1" go test -run TestClosures
func TestClosures(t *testing.T) {
	testInput := `
	let addxy = fn(x) {
		fn(y) {
			x + y;
		}
	}

	let addThem = addxy(22);
	addThem(20)
	`

	// Closures are functions that “close over” the environment they were defined in.
	// They carry their own environment around and whenever they’re called they can access it.

	testIntegerObject(t, testEval(testInput), 42)
}

// GOFLAGS="-count=1" go test -run TestStringLiteral
func TestStringLiteral(t *testing.T) {
	testInput := `"guten tag!"`

	evaluated := testEval(testInput)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "guten tag!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// GOFLAGS="-count=1" go test -run TestStringConcatenation
func TestStringConcatenation(t *testing.T) {
	testInput := `"guten" + " " + "tag!"`

	evaluated := testEval(testInput)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "guten tag!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// GOFLAGS="-count=1" go test -run TestBuiltinFunctions
func TestBuiltinFunctions(t *testing.T) {
	testInputs := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("guten tag!")`, 10},
		{`len(42)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		switch expected := ti.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

// GOFLAGS="-count=1" go test -run TestArrayLiterals
func TestArrayLiterals(t *testing.T) {
	testInput := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(testInput)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not an Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("Array len has wrong value. got=%d", len(arr.Elements))
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

// GOFLAGS="-count=1" go test -run TestArrayIndexExpressions
func TestArrayIndexExpressions(t *testing.T) {
	testInputs := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		switch expected := ti.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			testNullObject(t, evaluated)
		}
	}
}

// GOFLAGS="-count=1" go test -run TestHashLiterals
func TestHashLiterals(t *testing.T) {
	testInput := `
	let two = "two";
	{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
	}
	`

	evaluated := testEval(testInput)
	hashes, ok := evaluated.(*object.Hashes)
	if !ok {
		t.Fatalf("object is not a Hashes. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(hashes.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(hashes.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, found := hashes.Pairs[expectedKey]
		if !found {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

// GOFLAGS="-count=1" go test -run TestHashIndexExpressions
func TestHashIndexExpressions(t *testing.T) {
	testInputs := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 6}[false]`,
			6,
		},
	}

	for _, ti := range testInputs {
		evaluated := testEval(ti.input)
		switch expected := ti.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			testNullObject(t, evaluated)
		}
	}
}

// TODO: test in REPL
/*
let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(tail(arr), push(accumulated, f(first(arr))));
    }
  };

  iter(arr, []);
};

let map = fn(arr, f) { let iter = fn(arr, accumulated) { if (len(arr) == 0) { accumulated } else { iter(tail(arr), push(accumulated, f(first(arr)))); } iter(arr, []); };

let a = [1, 2, 3, 4];
let double = fn(x) { x * 2 };
map(a, double);
*/
