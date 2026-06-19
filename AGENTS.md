# Developer Instructions for RustX Starter

## Coding Conventions
- Use Go idiomatic style for the implementation.
- AST nodes should implement the `Node` interface.
- Evaluator should use tree-walking for now.

## Adding Built-ins
To add a new built-in function:
1. Open `evaluator/builtins.go`.
2. Add a new entry to the `builtins` map.
3. Implement the logic using `object.BuiltinFunction` signature.

## Scoping
- The language uses lexical scoping.
- `object.Environment` handles variable storage and nested scopes.
- `VarStatement` (`var` or `:=`) creates a new variable in the current scope.
- `AssignExpression` (`=`) updates an existing variable in the current or outer scope.

## Testing
Run all tests:
```bash
go test ./...
```
