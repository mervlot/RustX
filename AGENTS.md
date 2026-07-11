# RustX Developer Guide

This document provides instructions for developers working on the RustX Starter language.

## Design Philosophy
- Keep it simple and readable.
- Modularize everything.
- Avoid global state; use `Environment` for scoping.
- Every AST node must implement `GetToken()` for error reporting.

## Adding a New Built-in Function
1. Open `evaluator/builtins.go`.
2. Add a new entry to the `builtins` map.
3. Use `object.Object` for arguments and return types.
4. Use `newError` or `newErrorAt` for runtime errors.

## Modifying the Grammar
1. Update `token/token.go` if new tokens are needed.
2. Update `lexer/lexer.go` to recognize the new tokens.
3. Define new AST nodes in `ast/ast.go`.
4. Update `parser/parser.go` to parse the new constructs.
5. Update `evaluator/evaluator.go` to handle the new AST nodes.

## Scoping
Variables are looked up recursively in the environment. `env.Set()` always sets in the *current* scope (shadowing), while `env.Update()` (used for `=` assignment) searches up the chain to find and update an existing variable.

## Testing
Always run tests before submitting:
```bash
go test ./...
```
Also, run the sample programs to ensure no regressions in end-to-end flow.
