# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. It is implemented in Go.

## Features

- **Variables**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: Standard arithmetic, comparison, and logical operators.
- **Functions**: Parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Printing**: `println()`, `print()`.
- **Standard Library**: Built-in functions for strings, math, time, and filesystem basics.

## Installation

Ensure you have Go installed.

```bash
go mod tidy
go build -o rustx main.go
```

## Usage

### Run a Script

```bash
./rustx run samples/hello.rx
# or shorthand
./rustx samples/hello.rx
```

### Interactive REPL

```bash
./rustx repl
# or simply
./rustx
```

### Version

```bash
./rustx version
```

## Examples

### Factorial

```rust
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

## Architecture

- **Lexer**: Tokenizes the source code.
- **Parser**: Recursive descent + Pratt parser generating an AST.
- **AST**: Abstract Syntax Tree nodes.
- **Evaluator**: Tree-walking interpreter.
- **Object System**: Runtime representation of values.
- **Environment**: Scoping and variable management.
