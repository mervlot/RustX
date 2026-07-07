# Developer Documentation (AGENTS.md)

This document provides architectural details for developers working on the RustX Starter language.

## Design Philosophy

The goal is to create a clean, modular foundation that can be extended with features like a borrow checker, ownership system, and native compilation in the future.

## Project Structure

- `token/`: Defines the basic tokens of the language.
- `lexer/`: Implementation of the tokenizer.
- `ast/`: AST node definitions. Every node implements the `Node` interface.
- `parser/`: Recursive descent parser using Pratt parsing for expressions.
- `object/`: Runtime object system and `Environment` for scoping.
- `evaluator/`: Tree-walking interpreter logic.
- `repl/`: Interactive shell implementation.
- `std/`: (Future) Standard library modules.
- `samples/`: Example programs in `.rx`.

## Adding a Built-in Function

To add a new built-in function, modify the `builtins` map in `evaluator/evaluator.go`.

```go
"my_func": {
    Fn: func(args ...object.Object) object.Object {
        // Implementation
        return result
    },
},
```

## Internal Scoping

Scoping is managed by the `Environment` struct in `object/environment.go`. It supports nesting via the `outer` pointer. Variable updates (`Update` method) look up the chain until the variable is found, whereas `Set` always creates/overwrites in the current scope.

## Multiline REPL

The REPL in `repl/repl.go` uses a simple brace/parentheses balancing algorithm to decide when to evaluate multiline input.
