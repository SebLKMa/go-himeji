package parser

import (
	"fmt"
	"strconv"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

// For evaluation order - precedence
const (
	_ int = iota
	LOWEST
	EQUALS        // ==
	LESSERGREATER // > or <
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	CALL          // myFunction(X)
)

// Function types for associating to each specific token type
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  tk.Token // current token
	peekToken tk.Token // next token

	errors []string

	prefixParseFns map[tk.TokenType]prefixParseFn
	infixParseFns  map[tk.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Ensures curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[tk.TokenType]prefixParseFn)
	p.registerPrefix(tk.IDENT, p.parseIdentifier)
	p.registerPrefix(tk.INT, p.parseIntegerLiteral)

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
		//if stmt != nil { parseStatement never returns nil
		program.Statements = append(program.Statements, stmt)
		//}
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Move to next token if semi-colon, semi-colon is optional in an expression
	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	doPrefix := p.prefixParseFns[p.curToken.Type]
	if doPrefix == nil {
		fmt.Printf("Failed to find prefixParseFn for %T\n", p.curToken.Type)
		return nil
	}
	leftExpr := doPrefix()
	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	// Do not move to next token.
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	// Do not move to next token.
	return lit
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

func (p *Parser) registerPrefix(tokenType tk.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tk.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
