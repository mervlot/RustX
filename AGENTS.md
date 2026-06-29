# RustX Starter - Developer Guide

Welcome to the RustX Starter codebase. This project is a minimal, clean, and scalable foundation for a programming language inspired by Rust and Go.

## Architecture

The project follows a modular architecture:

- `token/`: Token definitions and lookup.
- `lexer/`: Lexical analysis (string to tokens).
- `ast/`: Abstract Syntax Tree node definitions.
- `parser/`: Recursive Descent + Pratt parser (tokens to AST).
- `object/`: Object-based runtime system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive REPL implementation.
- `samples/`: Example programs in RustX.

## Key Design Principles

1.  **Extensibility**: AST nodes and Objects are designed to be easily extended for future features like ownership, traits, and generics.
2.  **Scoping**: The `Environment` supports nested scoping. Variable updates (`=`) search up the scope chain, while initializations (`:=` or `var`) happen in the current scope.
3.  **Error Handling**: Precise line and column information is tracked from the lexer through to the evaluator for clear diagnostic messages.

## How to Add a New Built-in

To add a new built-in function:
1.  Open `evaluator/evaluator.go`.
2.  Add a new entry to the `builtins` map.
3.  Define the `Fn` using the `object.BuiltinFunction` signature.

## Internal Scoping

When implementing control flow (like `for` loops), remember to use `object.NewEnclosedEnvironment(env)` to ensure that variables defined within the loop do not leak into the outer scope.

## Testing

Run all tests from the root directory:
```bash
go test ./...
```
