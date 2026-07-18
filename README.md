# RustX Starter

A minimal, clean, scalable foundation for RustX, a compiled/interpreted hybrid programming language inspired by Rust and Go. It is focused on simplicity, speed, readability, and future extensibility.

Initially, RustX behaves like a lightweight scripting + systems hybrid language.

## Key Features

1. **Variables**: Declared with `var` keyword or using the shorthand assignment `:=`. Variables are mutable by default.
2. **Data Types**: Supports integers, floats, booleans, and strings.
3. **Operators**: Standard arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), and postfix increment/decrement (`++`, `--`).
4. **Control Flow**: Unparenthesized or parenthesized conditions in `if`, `else`, `while`, and Go-style `for` loops.
5. **Functions**: Supports parameters, return types, and recursion. Functions can be defined as named functions or anonymous function literals.
6. **Comments**: Both single-line (`//`) and multi-line (`/* */`) comments are skipped during lexing.
7. **Basic Standard Library**:
   - **Printing**: `print`, `println`
   - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
   - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
   - **Time**: `sleep` (milliseconds), `now` (Unix millisecond timestamp)
   - **Filesystem**: `read_file`, `write_file`
8. **Position-Aware Error Reporting**: Comprehensive line and column tracking for parser syntax errors and runtime evaluation errors.

## Code Structure

The project follows a modular, clean, and extensible architecture:
- `token/`: Defines lexical tokens and metadata.
- `lexer/`: Converts source code strings into a stream of tokens.
- `ast/`: Extensible Abstract Syntax Tree node definitions.
- `parser/`: Recursive Descent and Pratt parser for parsing statements and expressions.
- `object/`: Object representation and runtime environment (scopes, recursive bindings).
- `evaluator/`: Tree-walking interpreter evaluating AST nodes into runtime Objects.
- `repl/`: Interactive prompt supporting multiline input and environment state persistence.
- `errors/`: Custom structured Syntax and Runtime Error types.
- `samples/`: Example RustX program files.
- `main.go`: The CLI entry point.

## Build and Run

To run the project, ensure you have Go 1.21+ installed.

### Interactive REPL

Start the REPL to write RustX code dynamically:
```bash
go run main.go
```
or
```bash
go run main.go repl
```

To exit the REPL, type `exit` or `quit`.

### Running Scripts

You can execute a `.rx` file directly:
```bash
go run main.go run samples/hello.rx
```
Or use the shorthand command by passing the file path directly:
```bash
go run main.go samples/factorial.rx
```

### Running Tests

All packages have robust test suites. Run all tests with:
```bash
go test ./...
```
