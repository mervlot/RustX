# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Functions**: Parameters, return types (optional), and recursion.
- **Control Flow**: `if/else`, `while`, and `for` loops.
- **Built-in Functions**: `println`, `print`, `len`, `abs`, `pow`, `now`, `sleep`, `read_file`, `write_file`, and more.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).

## Installation

Ensure you have Go installed.

```bash
go build -o rustx main.go
```

## Usage

### Run a Script

```bash
./rustx run samples/hello.rx
```

### Interactive REPL

```bash
./rustx
```

### Version

```bash
./rustx version
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
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
println(factorial(5))
```

### Loops
```rust
for i := 0; i < 10; i++ {
    println(i)
}
```

## Architecture

The project follows a modular architecture:

- `/lexer`: Tokenizes the input string.
- `/token`: Defines the token types and structures.
- `/parser`: Builds an Abstract Syntax Tree (AST) using recursive descent and Pratt parsing.
- `/ast`: Defines the AST nodes.
- `/evaluator`: Executes the AST.
- `/object`: Defines the runtime object system and environment for scoping.
- `/repl`: Interactive Read-Eval-Print Loop.
- `main.go`: Entry point for the CLI.

## Future Extensibility

Designed with future additions in mind:
- Ownership and borrow checker
- Static typing and type checking
- Modules and imports
- Structs and Enums
- Concurrency (Goroutine-like)
- Bytecode VM or native compilation
