# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go. It is designed to be a foundation for a more advanced language, focused on simplicity, speed, and readability.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if/else`, `while`, and `for` loops (Go-style).
- **Operators**: Standard arithmetic and logical operators.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Built-in Functions**: Printing, math, string manipulation, time, and filesystem basics.

## Installation

Ensure you have [Go](https://golang.org/) installed.

```bash
git clone https://github.com/username/rustx.git
cd rustx
go build -o rustx main.go
```

## Usage

### Running a script

```bash
./rustx run samples/hello.rx
```

### Interactive REPL

```bash
./rustx
# OR
./rustx repl
```

### Version

```bash
./rustx version
```

## Syntax Examples

### Variables and Functions

```rust
func add(a int, b int) int {
    return a + b
}

var x = 10
y := 20
println(add(x, y))
```

### Control Flow

```rust
if x < y {
    println("x is smaller")
} else {
    println("x is larger or equal")
}

for i := 0; i < 5; i++ {
    println(i)
}
```

### Recursion

```rust
func factorial(n int) int {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
println(factorial(5))
```

## Project Architecture

- `lexer/`: Tokenizes the source code.
- `token/`: Defines the tokens used by the lexer and parser.
- `parser/`: Generates an Abstract Syntax Tree (AST) from tokens.
- `ast/`: Defines the AST nodes.
- `evaluator/`: Evaluates the AST nodes.
- `object/`: Defines the object system used during evaluation.
- `repl/`: Implements the interactive shell.
- `main.go`: Entry point for the CLI.

## License

MIT
