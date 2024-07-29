package object

import (
	"bytes"
	"fmt"
	"interpreter/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTIOn"
	STRING_OBJ   = "STRING"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
)

type BuiltinFunction func(args ...Object) Object

var _ Object = (*Builtin)(nil)

type Builtin struct {
	Fn BuiltinFunction
}

// Inspect implements Object.
func (b *Builtin) Inspect() string {
	return "builtin function "
}

// Type implements Object.
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

var _ Object = (*Integer)(nil)

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

var _ Object = (*Boolean)(nil)

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

var _ Object = (*Null)(nil)

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

var _ Object = (*ReturnObject)(nil)

type ReturnObject struct {
	Value Object
}

// Inspect implements Object.
func (r *ReturnObject) Inspect() string {
	return r.Value.Inspect()
}

// Type implements Object.
func (r *ReturnObject) Type() ObjectType {
	return RETURN_OBJ
}

var _ Object = (*Error)(nil)

type Error struct {
	Message string
}

// Inspect implements Object.
func (e *Error) Inspect() string {

	return "ERROR: " + e.Message
}

// Type implements Object.
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

var _ Object = (*FunctionObject)(nil)

type FunctionObject struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Inspect implements Object.
func (f *FunctionObject) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	return out.String()
}

// Type implements Object.
func (f *FunctionObject) Type() ObjectType {
	return FUNCTION_OBJ
}

var _ Object = (*StringObject)(nil)

type StringObject struct {
	Value string
}

// Inspect implements Object.
func (s *StringObject) Inspect() string {
	return s.Value
}

// Type implements Object.
func (s *StringObject) Type() ObjectType {
	return STRING_OBJ
}

var _ Object = (*ArrayObject)(nil)

type ArrayObject struct {
	Elements []Object
}

// Inspect implements Object.
func (a *ArrayObject) Inspect() string {
	var out bytes.Buffer
	elems := []string{}
	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

// Type implements Object.
func (a *ArrayObject) Type() ObjectType {
	return ARRAY_OBJ
}
