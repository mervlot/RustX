# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, focused on simplicity, speed, and readability.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Standard Library**: Built-in functions for math, strings, time, and filesystem.
- **CLI**: Execute scripts, start a REPL, or check the version.

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

### Start REPL

```bash
./rustx repl
```

### Check Version

```bash
./rustx version
```

## Syntax Example

```rustx
func factorial(n int) int {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

## Architecture

The project is modularized into several packages:
- `lexer`: Tokenizes the source code.
- `token`: Defines token types.
- `parser`: Recursive descent parser that generates the AST.
- `ast`: Abstract Syntax Tree nodes.
- `object`: Runtime object system.
- `evaluator`: Tree-walking interpreter.
- `repl`: Interactive REPL implementation.
