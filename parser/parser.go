package parser

import (
	"fmt"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

type Parser struct {
	l *lexer.Lexer

	curToken  tk.Token // current token
	peekToken tk.Token // next token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Ensures curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tk.TokenType) {
	msg := fmt.Sprintf("expected next token is %s, but got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != tk.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tk.LET:
		return p.parseLetStatement()
	case tk.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// Next token is expected to be an identifier
	// If so, move to next peek. Otherwise, fails
	if !p.moveNextIfPeekTokenIs(tk.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Next token is expected to be an equal sign
	// If so, move to next peek. Otherwise, fails
	if !p.moveNextIfPeekTokenIs(tk.ASSIGN) {
		return nil
	}

	// TODO: Ignoring expressions after the ASSIGN token, skip to semicolon for now
	for !p.curTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // move to next token after the "return" statement

	// TODO: Ignoring expressions after the ASSIGN token, skip to semicolon for now
	for !p.curTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t tk.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tk.TokenType) bool {
	return p.peekToken.Type == t
}

// moveNextIfPeekTokenIs moves current tokem to next token if the peek token is t.
// It allows the parser to assert correctness of the input statement.
func (p *Parser) moveNextIfPeekTokenIs(t tk.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken() // move to next token
		return true
	}
	p.peekError(t)
	return false
}
