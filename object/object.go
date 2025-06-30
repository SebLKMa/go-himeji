package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/seblkma/go-himeji/ast"
)

type ObjectType string

const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
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

type HashKey struct {
	Type  ObjectType
	Value uint64 // always a +ve hashed value
}

type Hashable interface {
	HashKey() HashKey
}

// A wrapper for integer with int64 value
type Integer struct {
	Value int64
}

// Implements the Object interface
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Implements the Object interface
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Implements Hashable
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// A wrapper for boolean with bool value
type Boolean struct {
	Value bool
}

// Implements the Object interface
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Implements the Object interface
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Implements Hashable
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

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
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewInnerEnvironment creates an inner Environment with a reference to its outer Environment
func NewInnerEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, found := e.store[name]
	// Mirrors variable scopes
	/*
		{
			a = 22
			{
				b = 20
				{
					c = a + b
				}
			}
		}
	*/
	if !found && e.outer != nil {
		obj, found = e.outer.Get(name)
	}
	return obj, found
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Implements Object interface
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// A wrapper for integer with string value
type String struct {
	Value string
}

// Implements the Object interface
func (s *String) Type() ObjectType { return STRING_OBJ }

// Implements the Object interface
func (s *String) Inspect() string { return s.Value }

// Implements Hashable
func (s *String) HashKey() HashKey {
	hFn := fnv.New64a()
	hFn.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: hFn.Sum64()}
}

type BuiltInFunction func(args ...Object) Object

// A wrapper for integer with string value
type Builtin struct {
	Fn BuiltInFunction
}

// Implements the Object interface
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Implements the Object interface
func (b *Builtin) Inspect() string { return "builtin function" }

// Our Array directly uses a Go slice
type Array struct {
	Elements []Object
}

// Implements the Object interface
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

// Implements the Object interface
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashes struct {
	Pairs map[HashKey]HashPair
}

// Implements the Object interface
func (h *Hashes) Type() ObjectType { return HASH_OBJ }

// Implements the Object interface
func (h *Hashes) Inspect() string {
	var out bytes.Buffer

	// We want to output the key value pair object associated to the key
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("[")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("]")

	return out.String()
}
