# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), Comparison (`==`, `!=`, `>`, `<`, `>=`, `<=`), Logical (`&&`, `||`).
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if-else`, `while`, and `for` loops.
- **Printing**: Built-in `print()` and `println()`.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Standard Library**: Basic support for strings, math, time, and filesystem.

## Getting Started

### Installation

Ensure you have Go installed on your system.

```bash
go build -o rustx main.go
```

### Running Scripts

You can run `.rx` files directly:

```bash
./rustx run samples/hello.rx
```

Or execute them without the `run` command:

```bash
./rustx samples/hello.rx
```

### Interactive REPL

Launch the REPL by running the executable without arguments:

```bash
./rustx
```

## Syntax Examples

### Variables

```rust
var x = 5
name := "RustX"
```

### Functions

```rust
func add(a int, b int) int {
    return a + b
}
```

### Loops

```rust
for (i := 0; i < 10; i++) {
    println(i)
}
```

## Architecture

The project follows a modular design:

- `/lexer`: Tokenizes source code.
- `/parser`: Recursive descent and Pratt parser generating an AST.
- `/ast`: Abstract Syntax Tree node definitions.
- `/object`: Object system and runtime environment.
- `/evaluator`: Tree-walking interpreter.
- `/repl`: Interactive shell.
- `/token`: Token definitions.

## Future Compatibility

Designed to eventually support:
- Ownership and Borrow Checker
- Async/Await and Concurrency
- Modules and Package Management
- Bytecode VM or Native Compilation
