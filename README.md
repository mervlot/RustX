# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), Comparison (`==`, `!=`, `>`, `<`, `>=`, `<=`), Logical (`&&`, `||`).
- **Functions**: Parameters, return types (optional), and recursion support.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Printing**: `print()` and `println()`.
- **Standard Library**: Basic support for strings, math, time, and filesystem.
- **Comments**: `//` single line and `/* */` multi-line.

## Usage

### Run a Script

```bash
go run main.go run samples/hello.rx
```

### Interactive REPL

```bash
go run main.go repl
```

### Version

```bash
go run main.go version
```

## Syntax Examples

### Hello World

```rustx
func main() {
    println("Hello World")
}
main()
```

### Variables

```rustx
var age = 16
name := "Muhammad"
```

### Factorial (Recursion)

```rustx
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### For Loop

```rustx
for (var i = 0; i < 10; i++) {
    println(i)
}
```

## Architecture

- **lexer/**: Tokenizer for RustX source code.
- **token/**: Token definitions.
- **ast/**: Abstract Syntax Tree nodes.
- **parser/**: Recursive descent + Pratt parser.
- **object/**: Runtime object system and environment.
- **evaluator/**: Tree-walking interpreter.
- **repl/**: Interactive shell implementation.

## Future Compatibility

The architecture is designed to support future additions like:
- Ownership and Borrow Checker
- Modules and Imports
- Concurrency (Goroutines/Channels)
- Structs, Enums, and Traits
- Bytecode VM or Native Compilation
