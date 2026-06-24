# RustX Starter Developer Guide

## Project Overview
RustX Starter is a foundation for a programming language inspired by Rust and Go. It currently features a tree-walking interpreter.

## Modular Architecture
- `token/`: Defines all supported tokens.
- `lexer/`: Converts source code into tokens.
- `ast/`: Contains the Abstract Syntax Tree node definitions.
- `parser/`: A recursive descent parser that also uses Pratt parsing for expressions.
- `object/`: The runtime object system and evaluation environment.
- `evaluator/`: Implements the tree-walking logic.
- `repl/`: Interactive shell.

## Adding a Built-in Function
To add a new built-in function, modify the `builtins` map in `evaluator/evaluator.go`.

## Scoping
- `Environment` handles variable scoping.
- Function calls create a new enclosed environment.
- For loops also use an enclosed environment for the loop header and body.

## Tips
- When adding new syntax, start with `token/token.go` and `lexer/lexer.go`.
- Ensure new AST nodes implement the `ast.Node` interface.
- Update the `Eval` function in `evaluator/evaluator.go` to handle new AST nodes.
- Use `samples/` to add regression tests/examples for new features.
