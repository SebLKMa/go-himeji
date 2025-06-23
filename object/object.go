package object

import "fmt"

type ObjectType string

const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

// The Object interface represents the internal representation of a value, e.g. integer, boolean, etc.
// Almost like "boxing" in C#
type Object interface {
	Type() ObjectType
	Inspect() string
}

// A wrapper for null, no value
type Null struct{}

// Implements the Object interface
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Implements the Object interface
func (n *Null) Inspect() string { return "null" }

// A wrapper for integer with int64 value
type Integer struct {
	Value int64
}

// Implements the Object interface
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Implements the Object interface
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// A wrapper for boolean with bool value
type Boolean struct {
	Value bool
}

// Implements the Object interface
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Implements the Object interface
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// A wrapper for boolean with bool value
type ReturnValue struct {
	Value Object
}

// Implements the Object interface
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Implements the Object interface
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Error struct {
	Message string
}

// Implements Object interface
func (e *Error) Type() ObjectType { return ERROR_OBJ }

func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// The Environment keeps track of objects bindings
type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, found := e.store[name]
	return obj, found
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
