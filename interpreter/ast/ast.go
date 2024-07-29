package ast

import (
	"bytes"
	"interpreter/token"
	"strings"
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

var _ Expression = (*InfixExpression)(nil)

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

// String implements Expression.
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// TokenLiteral implements Expression.
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *InfixExpression) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*Boolean)(nil)

type Boolean struct {
	Token token.Token
	Value bool
}

// String implements Expression.
func (b *Boolean) String() string {
	return b.Token.Literal
}

// TokenLiteral implements Expression.
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// expressionNode implements Expression.
func (b *Boolean) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*IfExpression)(nil)

type IfExpression struct {
	Token     token.Token
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

// String implements Expression.
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Then.String())
	if i.Else != nil {
		out.WriteString("else")
		out.WriteString(i.Else.String())
	}
	return out.String()

}

// TokenLiteral implements Expression.
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *IfExpression) expressionNode() {
	panic("unimplemented")
}

var _ Statement = (*BlockStatement)(nil)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

// String implements Statement.
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("\n}")
	return out.String()
}

// TokenLiteral implements Statement.
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

// statementNode implements Statement.
func (b *BlockStatement) statementNode() {
	panic("unimplemented")
}

var _ Expression = (*FunctionLiteral)(nil)

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

// String implements Expression.
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	return out.String()
}

// TokenLiteral implements Expression.
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// expressionNode implements Expression.
func (f *FunctionLiteral) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*CallExpression)(nil)

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

// String implements Expression.
func (c *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// TokenLiteral implements Expression.
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

// expressionNode implements Expression.
func (c *CallExpression) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*StringLiteral)(nil)

type StringLiteral struct {
	Token token.Token
	Value string
}

// String implements Expression.
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

// TokenLiteral implements Expression.
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

// expressionNode implements Expression.
func (s *StringLiteral) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*ArrayLiteral)(nil)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

// String implements Expression.
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer
	elems := []string{}
	for _, el := range a.Elements {
		elems = append(elems, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

// TokenLiteral implements Expression.
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

// expressionNode implements Expression.
func (a *ArrayLiteral) expressionNode() {
	panic("unimplemented")
}

var _ Expression = (*IndexExpression)(nil)

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

// String implements Expression.
func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

// TokenLiteral implements Expression.
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *IndexExpression) expressionNode() {
	panic("unimplemented")
}
