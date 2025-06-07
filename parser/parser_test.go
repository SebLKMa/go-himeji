package parser

import (
	"fmt"
	"testing"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
)

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
		t.Fatalf("Program statements expected %d, got %d\n", 3, len(program.Statements))
	}
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
		t.Errorf("s.TokenLiteral expected 'let' but got=%s\n", s.TokenLiteral())
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
