package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/username/rustx/ast"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	ARRAY_OBJ   = "ARRAY"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ERROR_OBJ    = "ERROR"
	NULL_OBJ     = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Primitive types
type Integer struct {
	Value int64
}
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

type Boolean struct {
	Value bool
}
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Null struct{}
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Composite / structural types
type Array struct {
	Elements []Object
}
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
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

// Environment represents scoped variables
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a variable in the current local block scope
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Update recursively searches outer scopes to update an existing variable.
// Returns false if the variable was not found in any scope.
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

// Functions & builtins
type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Environment
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		pStr := p.Name.String()
		if p.Type != "" {
			pStr += " " + p.Type
		}
		params = append(params, pStr)
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Error struct {
	Message string
	Line    int
	Column  int
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	return fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", e.Message, e.Line, e.Column)
}
