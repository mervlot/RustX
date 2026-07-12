# Agent Instructions for RustX Starter

This repository contains the implementation of RustX Starter, a programming language inspired by Rust and Go.

## Project Structure

- `token/`: Token definitions and lexer constants.
- `lexer/`: Lexical analyzer.
- `ast/`: Abstract Syntax Tree nodes.
- `parser/`: Recursive descent parser with Pratt parsing for expressions.
- `object/`: Runtime object system and environment management.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive REPL implementation.
- `std/`: (Integrated in evaluator/builtins.go) Standard library functions.
- `main.go`: CLI entry point.

## Adding Built-in Functions

To add a new built-in function:
1. Open `evaluator/builtins.go`.
2. Add a new entry to the `builtins` map.
3. Implement the logic using the `object.BuiltinFunction` signature.

## Internal Scoping

The interpreter uses an `Environment` object to manage variable scoping. Each block statement and function call creates a new enclosed environment to maintain proper lexical scope.

## Execution

- Run a file: `go run main.go run <file.rx>` or simply `go run main.go <file.rx>`.
- Start REPL: `go run main.go repl` or `go run main.go`.
- Check version: `go run main.go version`.

## Testing

Run all tests using:
```bash
go test ./...
```
