# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go. It is designed to be a lightweight scripting and systems hybrid language, with a focus on simplicity, speed, and readability.

## Features

- **Variables**: `var` keyword and `:=` shorthand for mutable variables.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Operators**: Standard arithmetic, comparison, and logical operators, including `++` and `--` postfix operators.
- **Functions**: Named and anonymous functions with parameter and return type support (types currently ignored by evaluator).
- **Control Flow**: `if`, `else`, `while`, and `for` loops.
- **Built-in Functions**: Comprehensive set of built-ins for math, strings, time, and filesystem operations.
- **REPL**: Interactive Read-Eval-Print Loop with multiline support.
- **Error Handling**: Human-readable syntax and runtime errors with line and column information.

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

### Loops
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

## Getting Started

### Installation
Ensure you have [Go](https://golang.org/) installed.

Clone the repository and build the project:
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
# or just
./rustx
```

#### Version
```bash
./rustx version
```

## Project Architecture

- `/lexer`: Lexical analyzer that converts source code into tokens.
- `/token`: Token definitions and lookup logic.
- `/ast`: Abstract Syntax Tree node definitions.
- `/parser`: Recursive descent + Pratt parser that builds the AST.
- `/object`: Runtime object system and environment management.
- `/evaluator`: Tree-walking interpreter that evaluates the AST.
- `/repl`: Interactive shell implementation.
- `/samples`: Collection of sample `.rx` programs.

## Future Compatibility
Designed as a foundation for:
- Ownership and Borrow Checker
- Modules and Package Management
- Concurrency (Goroutines/Channels)
- Structs, Enums, and Traits
- Bytecode VM or Native Compilation
