# RustX Developer Guide

## Core Concepts

RustX Starter is designed as a modular foundation for a modern, systems-leaning language.

### Packages

- `token`: Definitions of all language tokens.
- `lexer`: Lexical analysis.
- `ast`: Abstract Syntax Tree nodes.
- `parser`: Recursive descent + Pratt parser.
- `object`: Runtime object system and `Environment` management.
- `evaluator`: Tree-walking interpreter.
- `repl`: Interactive shell with multiline support.
- `main.go`: CLI entry point.

### Adding Built-ins

Built-ins are defined in `evaluator/evaluator.go` within the `builtins` map. To add a new built-in:
1. Define it in the `builtins` map with an `object.Builtin` wrapper.
2. Ensure it handles arguments and returns `object.Object` (or `object.Error` if something goes wrong).

### Scoping and Updates

The `Environment` in `object/object.go` supports:
- `Get`: Retrieve a variable, searching through outer scopes.
- `Set`: Create a variable in the current scope.
- `Update`: Modify an existing variable in its defining scope.

### Control Flow

Control flow statements (`if`, `while`, `for`) support optional parentheses. Semicolons in `for` loops are mandatory between clauses but optional after the loop itself.

## Syntax conventions

- `func` for functions.
- `var` for mutable variables.
- `:=` for shorthand assignment.
- `++` post-increment.
- `&&`, `||` logical operators (with short-circuiting).
- Optional semicolons.
