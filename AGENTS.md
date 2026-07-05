# RustX Starter Agent Instructions

You are Jules, a skilled software engineer working on RustX Starter.

## Project Overview

RustX Starter is a foundation-level programming language inspired by Rust and Go. It uses a modular architecture:
- `token`: Token definitions.
- `lexer`: Lexical analysis.
- `ast`: Abstract Syntax Tree.
- `parser`: Pratt Parser for expressions and recursive descent for statements.
- `object`: Object system and environments for the interpreter.
- `evaluator`: Tree-walking evaluator.
- `repl`: Interactive shell.

## Implementation Details

- **Module Path**: `github.com/username/rustx`
- **File Extension**: `.rx`
- **Execution**: `go run main.go [file.rx]` or `./rustx [file.rx]`
- **Testing**: `go test ./...` from the root.

## Adding Built-ins

To add a new built-in function, modify `evaluator/evaluator.go` and add it to the `builtins` map.

## Internal Scoping

Scoping is handled by the `Environment` struct in `object/object.go`. Use `NewEnclosedEnvironment` for nested scopes (functions, loops).

## Error Handling

Diagnostics should include line and column information from the `token.Token` associated with the `ast.Node`. Use `newError` in the evaluator for runtime errors.
