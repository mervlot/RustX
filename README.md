# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed as a foundation for a future systems language, currently implemented as a tree-walking interpreter in Go.

## Features

- **Simple Syntax**: Blends Rust and Go styles.
- **Variables**: Support for `var` and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Functions**: Parameters, return types, and recursion.
- **Control Flow**: `if/else`, `while`, and `for` loops.
- **REPL**: Interactive shell with multiline support.
- **Built-ins**: Standard library functions for math, strings, time, and filesystem.

## Installation

Ensure you have Go installed.

```bash
# Clone the repository (if applicable)
# go mod download
```

## Usage

### Run a Script

```bash
go run main.go run samples/hello.rx
```

### Start the REPL

```bash
go run main.go repl
# or just
go run main.go
```

### Version

```bash
go run main.go version
```

## Syntax Examples

### Hello World
```rust
func main() {
    println("Hello World")
}
main()
```

### Factorial (Recursion)
```rust
func factorial(n int) int {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}
println(factorial(5))
```

### Loops
```rust
for (i := 0; i < 10; i++) {
    println(i)
}
```

## Architecture

- `lexer/`: Tokenizes source code.
- `parser/`: Recursive descent + Pratt parser generating the AST.
- `ast/`: AST node definitions.
- `object/`: Runtime object system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive REPL.
- `main.go`: CLI entry point.
