# Development Tips for RustX Starter

## Code Organization
- **Lexer**: Add new tokens in `token/token.go` and update `lexer/lexer.go`.
- **Parser**: Uses Pratt parsing for expressions. Register new prefix/infix functions in `parser/parser.go`.
- **Evaluator**: Tree-walking style. Add new cases in `evaluator/evaluator.go`.

## Built-ins
To add a new built-in function:
1. Open `evaluator/evaluator.go`.
2. Add a new entry to the `builtins` map.

## Scoping
The `Environment` in `object/environment.go` handles nested scopes. Use `NewEnclosedEnvironment` for blocks and functions.

## Testing
Always run tests before submitting changes:
```bash
go test ./...
```
