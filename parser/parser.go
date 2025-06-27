package parser

import (
	"fmt"
	"strconv"

	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/lexer"
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

// For evaluation order precedence, the whole idea of PRATT parser
// the first LOWEST has lowest precedence, the last has the highest precedence
const (
	_ int = iota
	LOWEST
	EQUALS        // ==
	LESSERGREATER // > or <
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	CALL          // callFunction(X)
	INDEX         // array[index]
)

// Precedence table, e.g. multiplication has higher precedence than addition
// The whole idea of PRATT parser
var precedences = map[tk.TokenType]int{
	tk.EQ:       EQUALS,
	tk.NOT_EQ:   EQUALS,
	tk.LT:       LESSERGREATER,
	tk.GT:       LESSERGREATER,
	tk.PLUS:     SUM,
	tk.MINUS:    SUM,
	tk.SLASH:    PRODUCT,
	tk.ASTERISK: PRODUCT,
	tk.LPAREN:   CALL,
	tk.LBRACKET: INDEX,
}

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

	// prefix functions
	p.prefixParseFns = make(map[tk.TokenType]prefixParseFn)
	p.registerPrefix(tk.IDENT, p.parseIdentifier)
	p.registerPrefix(tk.INT, p.parseIntegerLiteral)
	p.registerPrefix(tk.BANG, p.parsePrefixExpression)
	p.registerPrefix(tk.MINUS, p.parsePrefixExpression)
	p.registerPrefix(tk.TRUE, p.parseBoolean)
	p.registerPrefix(tk.FALSE, p.parseBoolean)
	p.registerPrefix(tk.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(tk.IF, p.parseIfExpression)
	p.registerPrefix(tk.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(tk.STRING, p.parseStringLiteral)
	p.registerPrefix(tk.LBRACKET, p.parseArrayLiteral)

	// infix functions
	p.infixParseFns = make(map[tk.TokenType]infixParseFn)
	p.registerInfix(tk.PLUS, p.parseInfixExpression)
	p.registerInfix(tk.MINUS, p.parseInfixExpression)
	p.registerInfix(tk.SLASH, p.parseInfixExpression)
	p.registerInfix(tk.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tk.EQ, p.parseInfixExpression)
	p.registerInfix(tk.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tk.LT, p.parseInfixExpression)
	p.registerInfix(tk.GT, p.parseInfixExpression)
	p.registerInfix(tk.LPAREN, p.parseCallExpression)
	p.registerInfix(tk.LBRACKET, p.parseIndexExpression)

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
		return p.parseExpressionStatement() // parses prefix, infix as well
	}
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

	// TODO done: Ignoring expressions after the ASSIGN token, skip to semicolon for now
	//for !p.curTokenIs(tk.SEMICOLON) {
	//	p.nextToken()
	//}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // move to next token after the "return" statement

	// TODO done: Ignoring expressions after the ASSIGN token, skip to semicolon for now
	//for !p.curTokenIs(tk.SEMICOLON) {
	//	p.nextToken()
	//}

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Move to next token if semi-colon, semi-colon is optional in an expression
	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionList(end tk.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) { // e.g. end of array
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(tk.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.moveNextIfPeekTokenIs(end) {
		return nil
	}

	return list
}

// parseExpression parses prefixes by lookup tables - heart of the Pratt parser
func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	doPrefix := p.prefixParseFns[p.curToken.Type]
	if doPrefix == nil {
		//fmt.Printf("Failed to find prefixParseFn for %T\n", p.curToken.Type)
		p.noParseFnError(p.curToken.Type)
		return nil
	}
	leftExpr := doPrefix()

	// Attempts to find an infix with higher precedence by advancing the tokens
	for !p.peekTokenIs(tk.SEMICOLON) && precedence < p.peekPrecedence() {
		doInfix := p.infixParseFns[p.peekToken.Type]
		if doInfix == nil {
			p.noParseFnError(p.peekToken.Type)
			return leftExpr
		}

		p.nextToken()

		leftExpr = doInfix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expr := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken() // moves to next token, to parse the rhs expression

	expr.Right = p.parseExpression(PREFIX) // recursive call to parseExpression

	return expr
}

func (p *Parser) noParseFnError(t tk.TokenType) {
	msg := fmt.Sprintf("no parse function for %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()                              // moves to next token, to parse the rhs expression
	expr.Right = p.parseExpression(precedence) // recursive call to parseExpression, get back the rhs identifier

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	// Do not move to next token.
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
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

func (p *Parser) parseBoolean() ast.Expression {
	defer untrace(trace("parseBoolean"))
	b := &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(tk.TRUE)}
	return b
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("parseGroupExpression"))

	// When called here mean current token is LPAREN, move next token and parse again
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	// If after parsing and next token is not RPAREN, then this is not what we expect
	if !p.moveNextIfPeekTokenIs(tk.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("parseBlockStatement"))

	blk := &ast.BlockStatement{Token: p.curToken}
	blk.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(tk.RBRACE) && !p.curTokenIs(tk.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blk.Statements = append(blk.Statements, stmt)
		}
		p.nextToken()
	}

	return blk
}

func (p *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("parseIfExpression"))

	expr := &ast.IfExpression{Token: p.curToken}

	// If after parsing and next token is not LPAREN, then this is not what we expect
	if !p.moveNextIfPeekTokenIs(tk.LPAREN) {
		return nil
	}

	// Move to next token and parse statement again
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	// If after parsing and next token is not RPAREN, then this is not what we expect
	if !p.moveNextIfPeekTokenIs(tk.RPAREN) {
		return nil
	}

	// If after parsing and next token is not LBRACE, then this is not what we expect
	if !p.moveNextIfPeekTokenIs(tk.LBRACE) {
		return nil
	}

	expr.TrueBlock = p.parseBlockStatement()

	// Must handle the ELSE or RBRACE as well

	if p.peekTokenIs(tk.ELSE) {
		p.nextToken()
		if !p.moveNextIfPeekTokenIs(tk.LBRACE) {
			return nil
		}
		expr.FalseBlock = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(tk.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(tk.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.moveNextIfPeekTokenIs(tk.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fnl := &ast.FunctionLiteral{Token: p.curToken}

	if !p.moveNextIfPeekTokenIs(tk.LPAREN) {
		return nil
	}

	fnl.Parameters = p.parseFunctionParameters()

	if !p.moveNextIfPeekTokenIs(tk.LBRACE) {
		return nil
	}

	fnl.Body = p.parseBlockStatement()

	return fnl
}

// parseCallArguments has been refactored to parseCallExpression
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(tk.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(tk.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.moveNextIfPeekTokenIs(tk.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseCallExpression(callFunction ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.curToken, Function: callFunction}
	//expr.Arguments = p.parseCallArguments()
	expr.Arguments = p.parseExpressionList(tk.RPAREN) // same as parseCallArguments but accepts the ending token
	return expr
}

func (p *Parser) parseStringLiteral() ast.Expression {
	defer untrace(trace("parseStringLiteral"))
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	// Do not move to next token.
	return lit
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(tk.RBRACKET)
	return array
}

// parseIndexExpression typically parses array index expression
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseIndexExpression"))
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	expr.Index = p.parseExpression(LOWEST)

	if !p.moveNextIfPeekTokenIs(tk.RBRACKET) {
		return nil
	}

	// end of array detected
	return expr
}
