package parser

import (
	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

type Parser struct {
	l *lexer.Lexer

	curToken  tk.Token // current token
	peekToken tk.Token // next token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Ensures curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
