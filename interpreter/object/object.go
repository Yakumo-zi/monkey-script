package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN"
	ERROR_OBJ   = "ERROR"
)

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
