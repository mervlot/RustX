# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go.

## Features

- **Variables**: `var` and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Functions**: Named and anonymous functions, return values, and recursion.
- **Control Flow**: `if`, `else`, `while`, and `for` loops.
- **Standard Library**: Built-in functions for strings, math, time, and filesystem.
- **REPL**: Interactive shell with state preservation.
- **CLI**: Run scripts directly from the terminal.

## Syntax Examples

### Variables
```rustx
var age = 16
name := "Muhammad"
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
for i := 0; i < 5; i++ {
    println(i)
}
```

## Architecture

The project follows a modular design:
- `/lexer`: Tokenizes input.
- `/token`: Defines language tokens.
- `/ast`: Abstract Syntax Tree nodes.
- `/parser`: Recursive descent and Pratt parser.
- `/object`: Runtime value representation.
- `/evaluator`: AST-based interpreter.
- `/repl`: Interactive prompt.
- `/std`: Standard library (integrated as built-ins).

## Getting Started

### Installation
Ensure you have Go installed.

```bash
cd rustx
go build -o rustx main.go
```

### Running Scripts
```bash
./rustx run examples/hello.rx
```

### REPL
```bash
./rustx repl
```

### Version
```bash
./rustx version
```
