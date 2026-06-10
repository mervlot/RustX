# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` and `:=` shorthand
- **Data Types**: `int`, `float`, `bool`, `string`, `array`
- **Functions**: Parameters, return types, and recursion
- **Control Flow**: `if`, `else`, `while`, `for` loops
- **Operators**: Arithmetic, comparison, and logical
- **REPL**: Interactive shell with multiline support
- **Standard Library**: Basic math, strings, time, and filesystem operations

## Installation

Ensure you have Go installed.

```bash
go build -o rustx main.go
```

## Usage

### Run a script

```bash
./rustx run samples/hello.rx
```

### Start the REPL

```bash
./rustx repl
```

### Show version

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
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### Loops

```rustx
for (i := 0; i < 10; i++) {
    println(i)
}
```

## Architecture

The project is modularly structured:

- `/lexer`: Tokenizes the source code.
- `/parser`: Builds the Abstract Syntax Tree (AST) using a Pratt parser for expressions.
- `/ast`: Defines the structure of the language.
- `/object`: Implements the runtime object system and environment.
- `/evaluator`: Performs tree-walking evaluation of the AST.
- `/repl`: Provides an interactive shell.
- `/token`: Defines the language tokens.

## Future Compatibility

Designed to eventually support:
- Ownership and borrow checking
- Concurrency (async/await)
- Structs, Enums, and Traits
- Bytecode VM or native compilation
