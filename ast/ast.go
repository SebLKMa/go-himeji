package ast

import (
	"bytes"
	"strings"

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

type InfixExpression struct {
	Token    tk.Token // operator e.g. +, -, *, / , ==, !=
	Left     Expression
	Operator string
	Right    Expression
}

// Implements Expression
func (ie *InfixExpression) expressionNode() {}

// Implements Node
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// Implements Node
func (ie InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token tk.Token // token.BOOL
	Value bool
}

// Implements Expression
func (b *Boolean) expressionNode() {}

// Implements Node
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// Implements Node
func (b *Boolean) String() string { return b.Token.Literal }

type BlockStatement struct {
	Token      tk.Token // the "{" token
	Statements []Statement
}

// Implements Expression
func (bs *BlockStatement) expressionNode() {}

// Implements Node
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// Implements Node
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type IfExpression struct {
	Token      tk.Token // the "if" token
	Condition  Expression
	TrueBlock  *BlockStatement
	FalseBlock *BlockStatement
}

// Implements Expression
func (ife *IfExpression) expressionNode() {}

// Implements Node
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }

// Implements Node
func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString("  ")
	out.WriteString(ife.TrueBlock.String())
	if ife.FalseBlock != nil {
		out.WriteString("else ")
		out.WriteString(ife.FalseBlock.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      tk.Token // the "fn" token
	Parameters []*Identifier
	Body       *BlockStatement
}

// Implements Expression
func (fnl *FunctionLiteral) expressionNode() {}

// Implements Node
func (fnl *FunctionLiteral) TokenLiteral() string { return fnl.Token.Literal }

// Implements Node
func (fnl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fnl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fnl.TokenLiteral())
	out.WriteString(tk.LPAREN)
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(tk.RPAREN)
	out.WriteString(fnl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     tk.Token // the "(" token
	Function  Expression
	Arguments []Expression
}

// Implements Expression
func (ce *CallExpression) expressionNode() {}

// Implements Node
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// Implements Node
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString(tk.LPAREN)
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(tk.RPAREN)

	return out.String()
}
