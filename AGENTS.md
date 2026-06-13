# Developer Documentation

## Internal Architecture

RustX Starter follows a classic interpreter design:

1.  **Lexer (`lexer/`)**: Converts source code into a stream of tokens. Supports single-line (`//`) and multi-line (`/* */`) comments.
2.  **Parser (`parser/`)**: A recursive descent parser that uses Pratt parsing for expressions to handle operator precedence efficiently.
3.  **AST (`ast/`)**: Defines the structure of the program. Designed to be extensible for future features like ownership and types.
4.  **Object System (`object/`)**: Every runtime value is an `Object`. Scoping is managed by the `Environment`.
5.  **Evaluator (`evaluator/`)**: A tree-walking interpreter that executes the AST.

## Adding Built-in Functions

To add a new built-in function:
1. Open `evaluator/evaluator.go`.
2. Add a new entry to the `builtins` map.
3. Implement the logic using the `object.Object` interface.

## Scoping

The language uses lexical scoping. Variables are declared using `var` or `:=`. `Update` method in `Environment` ensures that assignments to existing variables correctly traverse up the scope chain, while `Declare` always creates a variable in the current scope.

## Future Extensibility

- **Ownership/Borrow Checker**: The AST nodes can be augmented with lifetime and ownership information. The evaluator or a separate pass can then enforce these rules.
- **Bytecode VM**: The AST can be compiled into bytecode instead of being interpreted directly.
- **Type Checking**: A separate semantic analysis pass can be added between parsing and evaluation.
