# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed to be simple, fast, and extensible, serving as a foundation for future features like ownership and borrow checking.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, and basic `Array`.
- **Operators**: Standard arithmetic, comparison, and logical operators.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, and `for` loops.
- **Built-in Functions**: Comprehensive set of string, math, time, and filesystem functions.
- **Interactive REPL**: Support for multiline input and persistent state.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)

### Installation

Clone the repository and build the CLI:

```bash
go build -o rustx main.go
```

### Usage

#### Run a Script

```bash
./rustx run samples/hello.rx
```

#### Start the REPL

```bash
./rustx repl
# Or simply
./rustx
```

#### Show Version

```bash
./rustx version
```

## Language Syntax

### Variable Declarations

```rustx
var x = 10
y := 20
x = 30 // Mutability
```

### Functions

```rustx
func add(a int, b int) int {
    return a + b
}

println(add(5, 10))
```

### Control Flow

```rustx
if (x > 10) {
    println("Greater")
} else {
    println("Smaller or equal")
}

for (i := 0; i < 5; i++) {
    println(i)
}

while (x > 0) {
    x = x - 1
}
```

## Architecture

The project is organized into modular components:

- `/token`: Token definitions and keyword mapping.
- `/lexer`: Lexical analyzer (tokenizer).
- `/ast`: Abstract Syntax Tree node definitions.
- `/parser`: Recursive descent parser with Pratt parsing for expressions.
- `/object`: Runtime object system and environment.
- `/evaluator`: Tree-walking interpreter and built-in functions.
- `/repl`: Interactive Read-Eval-Print Loop.
- `main.go`: CLI entry point.

## Testing

Run the full test suite:

```bash
go test ./...
```

## License

MIT
