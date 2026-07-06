# Developer Instructions for RustX Starter

This file contains tips and instructions for working on the RustX Starter codebase.

## Adding Built-ins

To add a new built-in function, modify the `builtins` map in `evaluator/evaluator.go`.
Ensure you validate the number and types of arguments to prevent panics.

## Scoping

The language uses a nested environment system. When updating a variable (e.g., `x = 10`),
it will look through the current environment and its ancestors to find the variable.
`var` and `:=` always create or update in the *current* environment (with `:=` allowing for shadowed variables in nested blocks if desired, though current implementation of `:=` in evaluator follows the same logic as `Set`).

*Update:* Actually, the `AssignExpression` in the evaluator distinguishes between `:=` (Set) and `=` (Update).

## Parsing

The parser uses a Pratt parsing approach for expressions. If you add new operators,
don't forget to:
1. Define the token in `token/token.go`.
2. Add the operator to the lexer in `lexer/lexer.go`.
3. Define its precedence in `parser/parser.go`.
4. Register the prefix or infix parse function in `parser/parser.go`.

## Error Handling

Always use `newError` in the evaluator to return an `object.Error`.
Use the constant format string to satisfy Go static analysis.
Syntax errors should be added to the parser's `errors` slice.
