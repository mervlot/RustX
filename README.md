# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `&&`, `||`, `++`, `--`.
- **Functions**: Parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Printing**: `println()`, `print()`.
- **Comments**: `//` single line, `/* */` multiline.
- **Standard Library**: Basic support for strings, math, time, and filesystem.

## Usage

### Run a script
```bash
go run main.go run samples/hello.rx
```

### Start REPL
```bash
go run main.go
```

### Version
```bash
go run main.go version
```

## Example

```rust
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

## Architecture

- `lexer/`: Lexical analyzer that converts source code into tokens.
- `token/`: Token definitions and types.
- `parser/`: Recursive descent parser that generates the AST.
- `ast/`: Abstract Syntax Tree node definitions.
- `object/`: Runtime object system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive Read-Eval-Print Loop.
- `main.go`: CLI entry point.
