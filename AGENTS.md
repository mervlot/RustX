# Developer Guide for RustX Starter

This document provides instructions for developers looking to extend or maintain RustX Starter.

## Adding a New Built-in Function

Built-in functions are registered in `evaluator/evaluator.go` within the `builtins` map.

1.  Open `evaluator/evaluator.go`.
2.  Add a new entry to the `builtins` map.
3.  Implement the logic using the `object.Object` interface.

Example:
```go
"sqrt": {
    Fn: func(args ...object.Object) object.Object {
        if len(args) != 1 {
            return newError("wrong number of arguments. got=%d, want=1", len(args))
        }
        if args[0].Type() != object.FLOAT_OBJ {
            return newError("argument to `sqrt` must be FLOAT, got %s", args[0].Type())
        }
        val := args[0].(*object.Float).Value
        return &object.Float{Value: math.Sqrt(val)}
    },
},
```

## Internal Scoping

Scoping is handled by the `Environment` struct in `object/object.go`.
- `Set`: Always sets the variable in the current scope.
- `Get`: Looks up a variable in the current scope, and if not found, recursively looks in the outer scope.
- `Update`: Looks up a variable and updates it where it was originally defined. This is used for assignment (`=`) and post-increment (`++`).

## Future Extensibility

The AST and Evaluator are designed to be extensible.
- To add a new statement or expression:
    1.  Add a new token type in `token/token.go`.
    2.  Add a new AST node in `ast/ast.go`.
    3.  Update the lexer in `lexer/lexer.go` to recognize the new syntax.
    4.  Update the parser in `parser/parser.go` to parse the new syntax into AST nodes.
    5.  Update the evaluator in `evaluator/evaluator.go` to handle the new AST node.

## Running Tests

Always run tests before submitting any changes:
```bash
go test ./...
```

To run a specific test:
```bash
go test -v ./lexer
```
