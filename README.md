# RustX Starter

RustX Starter is a minimal but clean, scalable programming language inspired by Rust and Go. It is designed as a foundation for a more advanced language, focused on simplicity, speed, and readability.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), Comparison (`==`, `!=`, `<`, `<=`, `>`, `>=`), Logical (`&&`, `||`, `!`), and Postfix (`++`, `--`).
- **Functions**: Parameters, return types (optional annotations), and recursion.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Built-in Functions**: `print`, `println`, `len`, `abs`, `pow`, `sin`, `cos`, `tan`, `contains`, `split`, `join`, `to_upper`, `to_lower`, `sleep`, `now`, `read_file`, `write_file`.
- **Comments**: Single-line (`//`) and multi-line (`/* */`).

## Installation

Ensure you have Go installed on your system.

```bash
go build -o rustx main.go
```

## Usage

### Run a script

```bash
./rustx run samples/hello.rx
```

Or using shorthand:

```bash
./rustx samples/hello.rx
```

### REPL

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
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}
println(factorial(5))
```

### Control Flow

```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

## Architecture

The project follows a modular architecture:

- `/token`: Token definitions.
- `/lexer`: Lexical analyzer (tokenizer).
- `/ast`: Abstract Syntax Tree nodes.
- `/parser`: Recursive descent + Pratt parser.
- `/object`: Runtime object system and environment.
- `/evaluator`: Tree-walking interpreter.
- `/repl`: Interactive Read-Eval-Print Loop.
- `/samples`: Sample programs.

## Future Extensibility

The architecture is designed to support future additions such as:
- Ownership and Borrow Checker
- Modules and Imports
- Concurrency (Goroutines/Channels)
- Structs, Enums, and Traits
- Bytecode VM or Native Compilation
