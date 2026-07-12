# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed with simplicity, speed, and readability in mind, serving as a foundation for future systems programming features like ownership and borrow checking.

## Features

- **Variable Bindings**: `var` for declarations and `:=` for shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, and `string`.
- **Functions**: Named and anonymous functions with recursion support.
- **Control Flow**: `if/else`, `while`, and `for` loops.
- **Standard Library**: Built-in support for math, strings, time, and filesystem operations.
- **Comments**: Single-line (`//`) and multi-line (`/* */`) comments.
- **REPL**: Interactive prompt with multiline input support.

## Syntax Examples

### Hello World
```rustx
func main() {
    println("Hello World")
}
```

### Factorial (Recursion)
```rustx
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### For Loop
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

## Project Architecture

- **Lexer**: Tokenizes source code.
- **Parser**: Recursive descent + Pratt parser for AST generation.
- **AST**: Position-aware Abstract Syntax Tree.
- **Object**: Runtime object system and environment.
- **Evaluator**: Tree-walking interpreter.
- **REPL**: Interactive shell with multiline support.

## Getting Started

### Installation

Ensure you have [Go](https://golang.org/) installed.

1. Clone the repository.
2. Initialize the module (if not already done): `go mod tidy`.

### Usage

**Run a script:**
```bash
go run main.go run samples/hello.rx
# or
go run main.go samples/hello.rx
```

**Interactive REPL:**
```bash
go run main.go repl
# or simply
go run main.go
```

**Check Version:**
```bash
go run main.go version
```

### Testing

Run all unit tests:
```bash
go test ./...
```

## Future Roadmap

- Ownership and Borrow Checker
- Modules and Imports
- Concurrency (Goroutine-style)
- Native Compilation / Bytecode VM
- Structs and Traits

## License

MIT
