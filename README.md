# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed to be simple, fast, and readable, with a focus on future extensibility.

## Features

- **Variables**: `var` and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `&&`, `||`, `++`, `--`.
- **Functions**: Parameters, return types, and recursion support.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Built-ins**: `print`, `println`, `len`, `push`, `abs`, `pow`, `to_upper`, `to_lower`, `split`, `now`, `sleep`, `read_file`, `write_file`.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).

## Installation

Ensure you have [Go](https://golang.org/) installed.

```bash
git clone <repository-url>
cd rustx
```

## Usage

### Run a script

```bash
go run main.go run samples/hello.rx
```

### REPL

```bash
go run main.go
```

### Version

```bash
go run main.go version
```

## Examples

### Hello World

```rust
func main() {
    println("Hello World")
}
main()
```

### Factorial

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

- **Lexer**: Tokenizes source code. Supports line/column tracking for error reporting.
- **Parser**: Recursive descent + Pratt parser. Generates an Abstract Syntax Tree (AST).
- **AST**: Extensible tree structure representing the program.
- **Evaluator**: Tree-walking interpreter that executes the AST.
- **Object System**: Internal representation of runtime values.
- **Environment**: Manages variable scoping and persistence.

## Future Plans

- Ownership and Borrow Checker
- Modules and Imports
- Concurrency (Goroutines style)
- Async/Await
- Structs and Enums
- Traits/Interfaces
- Bytecode VM or Native Compilation
