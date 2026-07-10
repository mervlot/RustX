# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed to be simple, fast, and readable, with a focus on future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Standard arithmetic, comparison, and logical operators.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if/else`, `while`, and `for` loops.
- **Built-in Functions**: `println`, `print`, `len`, `abs`, `pow`, `now`, `sleep`, `read_file`, `write_file`, and more.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **REPL**: Interactive Read-Eval-Print Loop.

## Installation

Ensure you have Go installed. Then, clone the repository and build the project:

```bash
go build -o rustx main.go
```

## Usage

### Running a script

```bash
./rustx run samples/hello.rx
```

Or directly:

```bash
./rustx samples/hello.rx
```

### Starting the REPL

```bash
./rustx repl
```

### Version

```bash
./rustx version
```

## Syntax Examples

### Hello World
```rustx
func main() {
    println("Hello World")
}
main()
```

### Factorial (Recursion)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}
println(factorial(5))
```

### Control Flow
```rustx
for i := 0; i < 10; i++ {
    if i % 2 == 0 {
        println(i, "is even")
    }
}
```

## Project Structure

- `lexer/`: Lexical analyzer.
- `parser/`: Recursive descent + Pratt parser.
- `ast/`: Abstract Syntax Tree definitions.
- `object/`: Runtime object system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive shell.
- `main.go`: CLI entry point.
