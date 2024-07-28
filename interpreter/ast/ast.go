package ast

import (
	"bytes"
	"interpreter/token"
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

var _ Node = &Program{}

type Program struct {
	Statements []Statement
}

// String implements Node.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

var _ Statement = (*LetStatement)(nil)

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// String implements Statement.
func (l *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// TokenLiteral implements Statement.
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// statementNode implements Statement.
func (l *LetStatement) statementNode() {
	panic("unimplemented")
}

var _ Expression = (*Identifier)(nil)

type Identifier struct {
	Token token.Token
	Value string
}

// String implements Expression.
func (i *Identifier) String() string {
	return i.Value
}

// TokenLiteral implements Expression.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *Identifier) expressionNode() {
	panic("unimplemented")
}

var _ Statement = (*ReturnStatement)(nil)

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// String implements Statement.
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// TokenLiteral implements Statement.
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// statementNode implements Statement.
func (r *ReturnStatement) statementNode() {
	panic("unimplemented")
}

var _ Statement = (*ExpressionStatement)(nil)

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// String implements Statement.
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// TokenLiteral implements Statement.
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// statementNode implements Statement.
func (e *ExpressionStatement) statementNode() {
	panic("unimplemented")
}

var _ Expression = (*IntegerLiteral)(nil)

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// String implements Expression.
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

// TokenLiteral implements Expression.
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *IntegerLiteral) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*PrefixExpression)(nil)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// String implements Expression.
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// TokenLiteral implements Expression.
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// expressionNode implements Expression.
func (p *PrefixExpression) expressionNode() {
	panic("unimplemented")
}
