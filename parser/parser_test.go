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
