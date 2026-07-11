# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed to be a foundation for a modern systems + scripting hybrid language, focused on simplicity, speed, and readability.

## Features

- **Blended Syntax**: Combines the best of Rust (`func`, curly braces) and Go (`:=`, simple types, loop style).
- **Hybrid Execution**: Runs as a script or via an interactive REPL.
- **Strong Foundation**: Modular architecture with Lexer, Parser (Recursive Descent + Pratt), AST, and Evaluator.
- **Built-in Types**: `int`, `float`, `bool`, `string`, `array`.
- **Control Flow**: `if/else`, `while`, and Go-style `for` loops.
- **Functions**: Support for parameters, return values, and recursion.
- **Standard Library**: Built-in functions for printing, math, time, and filesystem operations.

## Installation

Ensure you have Go (1.21+) installed.

```bash
# Clone the repository
git clone https://github.com/username/rustx
cd rustx

# Build the executable
go build -o rustx main.go
```

## Usage

### Run a Script
```bash
./rustx run samples/hello.rx
# or shorthand
./rustx samples/hello.rx
```

### Start REPL
```bash
./rustx
# or
./rustx repl
```

### Version
```bash
./rustx version
```

## Syntax Example

```rustx
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

## Project Architecture

- `/token`: Token definitions and lookup.
- `/lexer`: Lexical analysis (tokenization) with position tracking.
- `/ast`: Abstract Syntax Tree nodes, designed for extensibility.
- `/parser`: Recursive descent + Pratt parser.
- `/object`: Runtime object system and environment management.
- `/evaluator`: Tree-walking interpreter and built-in functions.
- `/repl`: Interactive prompt with multiline support.
- `main.go`: CLI entry point.

## Future Extensibility

RustX Starter is designed so that the following features can be added without a full rewrite:
- Ownership and Borrow Checker
- Advanced Type System (Structs, Enums, Traits)
- Concurrency (Goroutine-like or Async/Await)
- Modules and Package Management
- Bytecode VM or Native Compiler (LLVM)
