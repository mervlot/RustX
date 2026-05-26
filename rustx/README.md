# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, implemented in Go. It is designed to be minimal, clean, and extensible.

## Features

- **Variables**: `var` and `:=` shorthand.
- **Types**: `int`, `float`, `bool`, `string`.
- **Operators**: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `&&`, `||`.
- **Functions**: Parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, `for`.
- **Built-ins**: `println()`, `print()`, `len()`, `time()`, `sqrt()`.

## Installation

Ensure you have Go installed.

```bash
cd rustx
go build -o rustx
```

## Usage

### Run a script
```bash
./rustx run examples/demo.rx
```

### Interactive REPL
```bash
./rustx
```

### Version
```bash
./rustx version
```

## Syntax Example

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

## Architecture

- `lexer/`: Tokenizes the input string.
- `token/`: Defines valid language tokens.
- `parser/`: Recursive descent + Pratt parser that generates the AST.
- `ast/`: AST node definitions.
- `object/`: Runtime object system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive loop.
- `main.go`: CLI entry point.
