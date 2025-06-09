package ast

import (
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is a root Node holding all the statements
type Program struct {
	Statements []Statement
}

// Implements the Node interface
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token tk.Token // token.LET
	Name  *Identifier
	Value Expression
}

// Implements Statement
func (ls *LetStatement) statementNode() {}

// Implements Node
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token tk.Token // token.RETURN
	Value Expression
}

// Implements Statement
func (rs *ReturnStatement) statementNode() {}

// Implements Node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type Identifier struct {
	Token tk.Token // token.IDENT
	Value string
}

// Implements Expression
func (i *Identifier) expressionNode() {}

// Implements Node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
