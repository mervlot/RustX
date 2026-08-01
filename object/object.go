package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/username/rustx/ast"
)

type Type string

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	ARRAY_OBJ   = "ARRAY"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN_VALUE"
	ERROR_OBJ   = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ = "BUILTIN"
)

type Object interface {
	Type() Type
	Inspect() string
}

// Hashable interface for object keys if we ever need them
type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

// Integer Object
type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Float Object
type Float struct {
	Value float64
}

func (f *Float) Type() Type      { return FLOAT_OBJ }
func (f *Float) Inspect() string { return fmt.Sprintf("%g", f.Value) }

// Boolean Object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// String Object
type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Array Object
type Array struct {
	Elements []Object
}

func (ao *Array) Type() Type { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Null Object
type Null struct{}

func (n *Null) Type() Type      { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

// ReturnValue Object wrapper
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type      { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error Object wrapper for run-time and compile-time issues
type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Type() Type { return ERROR_OBJ }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError: %s at line %d column %d", e.Message, e.Line, e.Column)
	}
	return "RuntimeError: " + e.Message
}

// Function Object representation
type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		s := p.Name.String()
		if p.Type != "" {
			s += " " + p.Type
		}
		params = append(params, s)
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

// BuiltinFunction type signature
type BuiltinFunction func(line, col int, args ...Object) Object

// Builtin Object wrapper
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string { return "builtin function" }

// Environment holds variable bindings with support for recursive scope updating and nested scopes
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new root environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a nested scope
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves an identifier's bound value traversing parent scopes
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Set binds a new value in the local scope (e.g., var x = y or shorthand x := y)
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Update searches parent scopes to modify an existing variable binding (e.g., assignment x = y)
func (e *Environment) Update(name string, val Object) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return false
}
