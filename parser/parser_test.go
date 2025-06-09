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
		fmt.Println("no error")
		return
	}

	t.Errorf("Parser has %d errors\n", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
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
