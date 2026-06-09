# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go. It is designed to be simple, fast, and readable, with a foundation that allows for future features like ownership and borrow checking.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Types**: `int`, `float`, `bool`, `string`.
- **Operators**: Standard arithmetic (+, -, *, /), comparison (==, !=, <, >, <=, >=), and logical (&&, ||, !).
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Built-ins**:
    - `println(...)`, `print(...)`: Output to stdout.
    - `len(obj)`: Returns the length of a string or array.
    - `abs(num)`: Returns the absolute value.
    - `pow(base, exp)`: Returns base raised to exp.
    - `sin(rad)`: Returns the sine of radians.
    - `sleep(ms)`: Pauses execution for given milliseconds.
    - `now()`: Returns current Unix timestamp.
    - `split(str, sep)`: Splits string by separator into an array.
    - `read_file(path)`: Reads file content as a string.
    - `write_file(path, content)`: Writes string content to file.
- **Comments**: `//` single-line and `/* */` multi-line.

## Installation

Ensure you have Go installed on your system.

```bash
go build -o rustx main.go
```

## Usage

### Run a Script

```bash
./rustx run samples/hello.rx
```

### REPL

```bash
./rustx
```

### Version

```bash
./rustx version
```

## Syntax Examples

### Variables
```rustx
var x = 10
y := 20
```

### Functions
```rustx
func factorial(n: int) int {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}
```

### For Loop
```rustx
for (i := 0; i < 10; i++) {
    println(i)
}
```

## Architecture

- `lexer/`: Tokenizes the source code.
- `parser/`: Recursive descent parser that generates the AST.
- `ast/`: Defines the Abstract Syntax Tree nodes.
- `object/`: Defines the runtime object system and environment.
- `evaluator/`: Tree-walking interpreter.
- `repl/`: Interactive REPL implementation.
