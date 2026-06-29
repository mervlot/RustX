# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go. It is designed as a foundation for a powerful hybrid language, focusing on simplicity, speed, and readability.

## Features

- **Simple Syntax**: Blends Rust and Go styles.
- **Variables**: Support for `var` and `:=` shorthand assignment.
- **Data Types**: Integers, Floats, Booleans, Strings, and Arrays.
- **Control Flow**: `if/else`, `while` loops, and `for` loops.
- **Functions**: Parameters, return types, and recursion support.
- **Operators**: Full suite of arithmetic, comparison, and logical operators.
- **Standard Library**: Built-in support for math, strings, time, and filesystem operations.
- **REPL**: Interactive prompt with multiline support.

## Installation

Ensure you have Go installed on your system.

```bash
git clone https://github.com/username/rustx.git
cd rustx
```

## Usage

### Run a Script

```bash
go run main.go run samples/hello.rx
```

Alternatively, you can pass the file directly:

```bash
go run main.go samples/hello.rx
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

### Variables and Loops

```rustx
name := "RustX"
var count = 0

for i := 0; i < 5; i++ {
    println("Iteration:", i)
}
```

### Functions and Recursion

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

## Project Structure

- `lexer/`: Lexical analysis.
- `parser/`: Recursive descent parsing.
- `ast/`: AST node definitions.
- `evaluator/`: Tree-walking interpreter.
- `object/`: Runtime object system.
- `repl/`: Interactive shell.

## License

MIT
