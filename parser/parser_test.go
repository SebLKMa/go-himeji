package parser

import (
	"fmt"
	"testing"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		//fmt.Println("no error")
		return
	}

	t.Errorf("Parser has %d errors\n", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral expected 'let' but got=%q\n", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement) // type conversion, may or may not be LetStatement
	if !ok {
		t.Errorf("s expected LetStatement but got=%T\n", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value expected %s but got=%s\n", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() expected %s but got=%s\n", name, letStmt.Name.TokenLiteral())
		return false
	}

	fmt.Printf("testLetStatement\n%+v\n", letStmt)
	return true
}

// Helper function to test the right value of prefix expression
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	// Checks numeric
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	// Check string
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {

		t.Errorf("exp not *ast.Identifier. got=%T", expr)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

// Helper function to test the boolean expression
func testBooleanLiteral(t *testing.T, il ast.Expression, value bool) bool {
	bl, ok := il.(*ast.Boolean)
	if !ok {
		t.Errorf("il not *ast.Boolean. got=%T", il)
		return false
	}

	// Checks numeric
	if bl.Value != value {
		t.Errorf("bl.Value not %t. got=%t", value, bl.Value)
		return false
	}
	// Check string
	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bl.TokenLiteral not %t. got=%s", value, bl.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of exp not handled. got=%T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", expr, expr)
		return false
	}
	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}
	if opExpr.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExpr.Operator)
		return false
	}
	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}
	return true
}

// GOFLAGS="-count=1" go test -run TestLetStatements
func TestLetStatements(t *testing.T) {
	input := `
	let x = 20;
	let y = 22;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program statements expected %d, but got %d\n", 3, len(program.Statements))
	}

	checkParserErrors(t, p)

	// Only tests the variable name, not the assignment.
	expectedIdentifiers := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, ei := range expectedIdentifiers {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, ei.expectedIdentifier) {
			return
		}
	}
}

// GOFLAGS="-count=1" go test -run TestLetStatementsError
func TestLetStatementsError(t *testing.T) {
	input := `
	let x 20;
	let = 22;
	let 838383;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program statements expected %d, but got %d\n", 3, len(program.Statements))
	}

	checkParserErrors(t, p)

	expectedIdentifiers := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, ei := range expectedIdentifiers {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, ei.expectedIdentifier) {
			return
		}
	}
}

// GOFLAGS="-count=1" go test -run TestLetStatementsComplete
func TestLetStatementsComplete(t *testing.T) {
	lets := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	// Tests both the variable name and the assignment.
	for _, let := range lets {
		l := lexer.New(let.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, let.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, let.expectedValue) {
			return
		}
	}

}

// GOFLAGS="-count=1" go test -run TestReturnStatements
func TestReturnStatements(t *testing.T) {
	input := `
	return 20;
	return 22;
	return 838383;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program statements expected %d, but got %d\n", 3, len(program.Statements))
	}

	for _, stmt := range program.Statements {
		rs, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt expected type is ast.ReturnStatement, but got %T\n", stmt)
			continue
		}
		if rs.TokenLiteral() != "return" {
			t.Errorf("stmt unexpected literal value %s\n", rs.TokenLiteral())
		}
	}
}

// GOFLAGS="-count=1" go test -run TestIdentifierExpression
func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("statement expected type is ast.Identifier, but got %T\n", stmt.Expression)
	}

	// Identifiers will not have semi-colon

	if ident.Value != "foobar" {
		t.Errorf("identifier expected value is foobar, but got %s\n", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("identifier unexpected literal value %s\n", ident.TokenLiteral())
	}

}

// GOFLAGS="-count=1" go test -run TestIntegerLiteralExpression
func TestIntegerLiteralExpression(t *testing.T) {
	input := `42;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.IntegerLiteral, but got %T\n", stmt.Expression)
	}

	// Identifiers will not have semi-colon

	if literal.Value != 42 {
		t.Errorf("identifier expected value is 42, but got %d\n", &literal.Value)
	}
	if literal.TokenLiteral() != "42" {
		t.Errorf("identifier unexpected literal value %s\n", literal.TokenLiteral())
	}

}

// GOFLAGS="-count=1" go test -run TestParsingPrefixExpressions
func TestParsingPrefixExpressions(t *testing.T) {
	prefixInputs := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, pi := range prefixInputs {
		l := lexer.New(pi.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if expr.Operator != pi.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", pi.operator, expr.Operator)
		}

		// Using testLiteralExpression instead
		/*
			if !testIntegerLiteral(t, expr.Right, pi.value) {
				return
			}
		*/
		if !testLiteralExpression(t, expr.Right, pi.value) {
			return
		}
	}
}

// GOFLAGS="-count=1" go test -run TestParsingInfixExpressions
func TestParsingInfixExpressions(t *testing.T) {
	/* input:
	5 + 5;
	5 - 5;
	5 * 5;
	5 / 5;
	5 > 5;
	5 < 5;
	5 == 5;
	5 != 5;
	*/
	infixInputs := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, ii := range infixInputs {
		l := lexer.New(ii.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, expr.Left, ii.leftValue) {
			t.Fatalf("expr.Left inconsistent with %d\n", ii.leftValue)
			return
		}

		if expr.Operator != ii.operator {
			t.Fatalf("expr.Operator is not '%s'. got=%s", ii.operator, expr.Operator)
			return
		}

		if !testLiteralExpression(t, expr.Right, ii.rightValue) {
			t.Fatalf("expr.Right inconsistent with %d\n", ii.rightValue)
		}

	}
}

// GOFLAGS="-count=1" go test -run TestOperatorPrecedenceParsing
func TestOperatorPrecedenceParsing(t *testing.T) {
	infixInputs := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		// adding boolean expressions
		{"true", "true"},
		{"false", "false"},
		{"3 < 5 == false", "((3 < 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		// adding group expressions
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		// adding call functions having highest precedence
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		// array index expressions
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, ii := range infixInputs {
		l := lexer.New(ii.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != ii.expected {
			t.Errorf("expected=%q, got=%q", ii.expected, actual)
		}
	}
}

// GOFLAGS="-count=1" go test -run TestInfixStatements
func TestInfixStatements(t *testing.T) {
	infixInput := "5 + 10"

	l := lexer.New(infixInput)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testInfixExpression(t, stmt.Expression, 5, "+", 10)

	infixInput = "alice * bob"

	l = lexer.New(infixInput)
	p = New(l)
	program = p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok = program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testInfixExpression(t, stmt.Expression, "alice", "*", "bob")
}

// GOFLAGS="-count=1" go test -run TestBooleanExpression
func TestBooleanExpression(t *testing.T) {
	/*
		true;
		false;
		let foobar = true;
		let barfoo = false;
	*/
	input := `true;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	b, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("statement expected type is ast.Boolean, but got %T\n", stmt.Expression)
	}

	// Identifiers will not have semi-colon

	if !b.Value {
		t.Errorf("identifier expected value is 42, but got %v\n", &b.Value)
	}
	if b.TokenLiteral() != "true" {
		t.Errorf("identifier unexpected literal value %s\n", b.TokenLiteral())
	}

}

// GOFLAGS="-count=1" go test -run TestIfExpression
func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement expected type is ast.IfExpression, but got %T\n", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.TrueBlock.Statements) != 1 {
		t.Fatalf("TrueBlock statements expected %d, but got %d\n", 1, len(expr.TrueBlock.Statements))
	}

	trueBlock, ok := expr.TrueBlock.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.TrueBlock.Statements expected type is ast.ExpressionStatement, but got %T\n", expr.TrueBlock.Statements[0])
	}

	if !testIdentifier(t, trueBlock.Expression, "x") {
		return
	}

	if expr.FalseBlock != nil {
		t.Errorf("expr.FalseBlock.Statements was not nil. got=%+v", expr.FalseBlock)
	}
}

// GOFLAGS="-count=1" go test -run TestIfElseExpression
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement expected type is ast.IfExpression, but got %T\n", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.TrueBlock.Statements) != 1 {
		t.Fatalf("TrueBlock statements expected %d, but got %d\n", 1, len(expr.TrueBlock.Statements))
	}

	trueBlock, ok := expr.TrueBlock.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.TrueBlock.Statements expected type is ast.ExpressionStatement, but got %T\n", expr.TrueBlock.Statements[0])
	}

	if !testIdentifier(t, trueBlock.Expression, "x") {
		return
	}

	if len(expr.TrueBlock.Statements) != 1 {
		t.Fatalf("TrueBlock statements expected %d, but got %d\n", 1, len(expr.TrueBlock.Statements))
	}

	falseBlock, ok := expr.FalseBlock.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.FalseBlock.Statements expected type is ast.ExpressionStatement, but got %T\n", expr.TrueBlock.Statements[0])
	}

	if !testIdentifier(t, falseBlock.Expression, "y") {
		return
	}

}

// GOFLAGS="-count=1" go test -run TestFunctionLiteral
func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myFunc, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.FunctionLiteral, but got %T\n", stmt.Expression)
	}

	if len(myFunc.Parameters) != 2 {
		t.Fatalf("My function parameters expected %d, but got %d\n", 2, len(myFunc.Parameters))
	}

	if !testLiteralExpression(t, myFunc.Parameters[0], "x") {
		t.Fatal("function parameter expected x but got ", myFunc.Parameters[0])
	}
	if !testLiteralExpression(t, myFunc.Parameters[1], "y") {
		t.Fatal("function parameter expected y but got ", myFunc.Parameters[1])
	}

	if len(myFunc.Body.Statements) != 1 {
		t.Fatalf("My function body statement expected %d, but got %d\n", 1, len(myFunc.Body.Statements))
	}

	bodyStmt, ok := myFunc.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("My function body statement expectedtype is ast.ExpressionStatement, but got %T\n", myFunc.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// GOFLAGS="-count=1" go test -run TestCallExpression
func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myCall, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement expected type is ast.CallExpression, but got %T\n", stmt.Expression)
	}

	if !testIdentifier(t, myCall.Function, "add") {
		return
	}

	if len(myCall.Arguments) != 3 {
		t.Fatalf("My call arguments expected %d, but got %d\n", 3, len(myCall.Arguments))
	}

	if !testLiteralExpression(t, myCall.Arguments[0], 1) {
		t.Fatal("call argument expected 1 but got ", myCall.Arguments[0])
	}
	if !testInfixExpression(t, myCall.Arguments[1], 2, "*", 3) {
		t.Fatal("call argument expression expected 2 * 3 but got ", myCall.Arguments[1])
	}
	if !testInfixExpression(t, myCall.Arguments[2], 4, "+", 5) {
		t.Fatal("call argument expression expected 4 + 5 but got ", myCall.Arguments[1])
	}
}

// GOFLAGS="-count=1" go test -run TestStringLiteralExpression
func TestStringLiteralExpression(t *testing.T) {
	input := `"guten tag!";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.StringLiteral, but got %T\n", stmt.Expression)
	}

	if literal.Value != "guten tag!" {
		t.Errorf("identifier expected value is guten tag!, but got %d\n", &literal.Value)
	}
}

// GOFLAGS="-count=1" go test -run TestArrayLiterals
func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myArray, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.ArrayLiteral, but got %T\n", stmt.Expression)
	}

	if len(myArray.Elements) != 3 {
		t.Fatalf("My array elements expected %d, but got %d\n", 3, len(myArray.Elements))
	}

	testInfixExpression(t, myArray.Elements[1], 2, "*", 2)
	testInfixExpression(t, myArray.Elements[2], 3, "+", 3)
}

// GOFLAGS="-count=1" go test -run TestIndexExpression
func TestIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	indexExpr, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("statement expected type is ast.IndexExpression, but got %T\n", stmt.Expression)
	}

	if !testIdentifier(t, indexExpr.Left, "myArray") {
		return
		//t.Fatalf("My array identifier expected %s, but got %s\n", "myArray", indexExpr.Left)
	}

	if !testInfixExpression(t, indexExpr.Index, 1, "+", 1) {
		return
	}

	// see TestOperatorPrecedenceParsing updated to also test precedence in array index expression
}

// GOFLAGS="-count=1" go test -run TestHashLiterals
func TestHashLiterals(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myHash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.HashLiteral, but got %T\n", stmt.Expression)
	}

	if len(myHash.Pairs) != 3 {
		t.Fatalf("My hash pairs expected %d, but got %d\n", 3, len(myHash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range myHash.Pairs {
		keyLiteral, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[keyLiteral.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

// GOFLAGS="-count=1" go test -run TestEmptyHashLiteral
func TestEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myHash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.HashLiteral, but got %T\n", stmt.Expression)
	}

	if len(myHash.Pairs) != 0 {
		t.Fatalf("My hash pairs expected %d, but got %d\n", 0, len(myHash.Pairs))
	}
}

// GOFLAGS="-count=1" go test -run TestParseHashLiteralsExpressions
func TestParseHashLiteralsExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements expected %d, but got %d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected type is ast.ExpressionStatement, but got %T\n", program.Statements[0])
	}

	myHash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("statement expected type is ast.HashLiteral, but got %T\n", stmt.Expression)
	}

	if len(myHash.Pairs) != 3 {
		t.Fatalf("My hash pairs expected %d, but got %d\n", 3, len(myHash.Pairs))
	}

	testFns := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range myHash.Pairs {
		keyLiteral, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		testFn, found := testFns[keyLiteral.String()]
		if !found {
			t.Errorf("No test function for key %q found", keyLiteral.String())
			continue
		}

		testFn(value)
	}
}
