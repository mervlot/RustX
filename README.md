# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed for simplicity, speed, and readability, providing a solid foundation for future features like an ownership system and borrow checker.

## Features

- **Variables**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Standard arithmetic and logical operators.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, and `for` loops.
- **Comments**: Single-line (`//`) and multi-line (`/* */`) support.
- **Standard Library**: Basic built-ins for printing, math, time, and filesystem.

## Installation

Ensure you have [Go](https://golang.org/) installed on your system.

```bash
go mod tidy
```

## Usage

### Run a Script
Execute a `.rx` file directly from the terminal:
```bash
go run main.go run samples/hello.rx
```

### Interactive REPL
Start the interactive REPL:
```bash
go run main.go
```

### Check Version
```bash
go run main.go version
```

## Syntax Examples

### Variables
```rustx
var x = 5
y := 10
```

### Functions
```rustx
func add(a, b) {
    return a + b
}

println(add(5, 10))
```

### Control Flow
```rustx
if x < y {
    println("x is less than y")
} else {
    println("x is not less than y")
}

for i := 0; i < 5; i++ {
    println(i)
}
```

## Project Structure

- `lexer/`: Lexical analysis and tokenization.
- `parser/`: Recursive descent and Pratt parsing.
- `ast/`: Abstract Syntax Tree node definitions.
- `evaluator/`: Tree-walking interpreter and built-ins.
- `object/`: Runtime value representation and environment.
- `repl/`: Interactive read-eval-print loop.
- `token/`: Token definitions.
- `samples/`: Example RustX programs.

## Testing

Run the full test suite:
```bash
go test ./...
```
