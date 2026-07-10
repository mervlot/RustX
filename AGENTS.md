# Developer Guide (AGENTS.md)

This document provides internal details about the RustX Starter implementation for developers and agents.

## Architecture

RustX Starter follows a classic interpreter architecture:
1. **Lexer** (`lexer/`): Converts source code into tokens.
2. **Parser** (`parser/`): Converts tokens into an Abstract Syntax Tree (AST). Uses Recursive Descent for statements and Pratt parsing for expressions.
3. **AST** (`ast/`): Defines the nodes of the tree.
4. **Object System** (`object/`): Defines the types used during evaluation (Integer, Boolean, Function, etc.).
5. **Evaluator** (`evaluator/`): Walks the AST and executes the logic.

## Key Design Decisions

- **Optional Semicolons**: The parser handles semicolons if present but does not strictly require them for most statements.
- **Unparenthesized Control Flow**: `if`, `while`, and `for` headers do not require parentheses around conditions, similar to Rust and Go.
- **Environment**: Scoping is handled by the `Environment` struct, which supports nested scopes via an `outer` pointer.
- **Built-ins**: Standard library functions are registered in `evaluator/evaluator.go`.

## Adding a New Built-in

To add a new built-in function:
1. Open `evaluator/evaluator.go`.
2. Find the `builtins` map.
3. Add a new entry with the function name and its implementation.

## Adding a New Token/Syntax

1. Add the token type to `token/token.go`.
2. Update the lexer in `lexer/lexer.go` to recognize the new syntax.
3. Add AST nodes in `ast/ast.go`.
4. Update the parser in `parser/parser.go` to parse the new syntax into AST nodes.
5. Update the evaluator in `evaluator/evaluator.go` to handle the new AST nodes.

## Testing

Always run tests before submitting:
```bash
go test ./...
```
Also verify changes with sample programs in `samples/`.
