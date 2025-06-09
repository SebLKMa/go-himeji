package ast

import (
	"bytes"

	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

type Node interface {
	TokenLiteral() string
	String() string
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

// Implements the Note interface
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// Implements Node
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token tk.Token // token.RETURN
	Value Expression
}

// Implements Statement
func (rs *ReturnStatement) statementNode() {}

// Implements Node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Implements Node
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(" = ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      tk.Token // the first token of the expression
	Expression Expression
}

// Implements Statement
func (es *ExpressionStatement) statementNode() {}

// Implements Node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// Implements Node
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token tk.Token // token.IDENT
	Value string
}

// Implements Expression
func (i *Identifier) expressionNode() {}

// Implements Node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Implements Node
func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
	Token tk.Token // token.INT
	Value int64
}

// Implements Expression
func (il *IntegerLiteral) expressionNode() {}

// Implements Node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// Implements Node
func (il *IntegerLiteral) String() string { return il.Token.Literal } // Value is int64

type PrefixExpression struct {
	Token    tk.Token // e.g. !, +, -, ...
	Operator string
	Right    Expression
}

// Implements Expression
func (pe *PrefixExpression) expressionNode() {}

// Implements Node
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// Implements Node
func (pe PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
