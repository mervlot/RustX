# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, implemented in Go. It is designed for simplicity, readability, and future extensibility.

## Features

- **Variables:** `var` keyword and `:=` shorthand.
- **Data Types:** `int`, `float`, `bool`, `string`.
- **Functions:** Support for parameters, return types, and recursion.
- **Control Flow:** `if/else`, `while`, and Go-style `for` loops.
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), and Logical (`&&`, `||`).
- **Comments:** Single-line (`//`) and multi-line (`/* */`) support.
- **REPL:** Interactive prompt for testing code snippets.

## Syntax Examples

### Variable Declaration
```rustx
var x = 5
name := "RustX"
```

### Functions
```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
```

### Loops
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

## Getting Started

### Installation
Ensure you have [Go](https://golang.org/dl/) installed.

1. Clone the repository.
2. Navigate to the `rustx` directory.
3. Build the CLI:
   ```bash
   go build -o rustx main.go
   ```

### Running Scripts
```bash
./rustx run examples/test.rx
```

### Interactive REPL
```bash
./rustx repl
```

### Version
```bash
./rustx version
```

## Architecture

- `lexer/`: Lexical analysis (tokenization).
- `parser/`: Recursive descent + Pratt parsing to generate the AST.
- `ast/`: Abstract Syntax Tree node definitions.
- `evaluator/`: Tree-walking interpreter.
- `object/`: Runtime object system and environment management.
- `repl/`: Interactive loop implementation.

## Testing
Run the test suite using:
```bash
go test ./...
```
