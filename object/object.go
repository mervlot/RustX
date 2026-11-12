package object

import (
	"bytes"
	"fmt"
	"rustx/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	FLOAT_OBJ        ObjectType = "FLOAT"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	STRING_OBJ       ObjectType = "STRING"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
)

// Object is an interface representing any runtime value.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer object
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Float object
type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

// Boolean object
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// String object
type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

// Null object
type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

// ReturnValue object wraps returned value to unwind call stack
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Error object
type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
	}
	return fmt.Sprintf("RuntimeError:\n%s", e.Message)
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Function object
type Function struct {
	Name       string
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.Name.Value
	}
	name := f.Name
	if name == "" {
		name = "anonymous"
	}
	out.WriteString("func " + name + "(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ... }")
	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// BuiltinFunction represents standard library builtins
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Array object
type Array struct {
	Elements []Object
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, len(ao.Elements))
	for i, e := range ao.Elements {
		elements[i] = e.Inspect()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

// Environment manages variable bindings and scoping
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new top-level environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// NewEnclosedEnvironment creates a new enclosed (child) environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get looks up a variable in the current or outer environments
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set defines a variable in the CURRENT environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Update recursively finds the environment where the variable exists and updates it.
// If it does not exist anywhere, updates current scope.
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
