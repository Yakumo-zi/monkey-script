package ast

import "interpreter/token"

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

var _ Node = &Program{}

type Program struct {
	Statements []Statement
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

// TokenLiteral implements Expression.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *Identifier) expressionNode() {
	panic("unimplemented")
}
