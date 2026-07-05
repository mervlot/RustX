# RustX Starter

RustX Starter is a clean, scalable programming language inspired by Rust and Go. It is designed to be simple, fast, and extensible, serving as a foundation for future systems programming features like ownership and borrow checking.

## Features

- **Variables**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Functions**: Support for parameters, return types, and recursion.
- **Control Flow**: `if/else`, `while`, and `for` loops.
- **Built-in Functions**: Comprehensive standard library for strings, math, time, and filesystem operations.
- **REPL**: Interactive shell with multiline support.
- **Clean Architecture**: Modular design (Lexer, Parser, AST, Evaluator, Object System).

## Installation

Ensure you have Go installed on your system.

```bash
git clone <repository-url>
cd rustx
go build -o rustx main.go
```

## Usage

### Run a script
```bash
./rustx run samples/hello.rx
# or directly
./rustx samples/hello.rx
```

### Start the REPL
```bash
./rustx repl
# or just
./rustx
```

### Show version
```bash
./rustx version
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

println(add(5, 10))
```

### Loops
```rust
for i := 0; i < 10; i++ {
    println(i)
}

var x = 5
while x > 0 {
    println(x)
    x--
}
```

### Conditionals
```rust
if name == "RustX" {
    println("Hello, RustX!")
} else {
    println("Hello, World!")
}
```

## Standard Library

- **Printing**: `print`, `println`
- **Strings**: `len`, `split`, `join`, `to_upper`, `to_lower`, `contains`
- **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
- **Time**: `now`, `sleep`
- **Filesystem**: `read_file`, `write_file`
- **Arrays**: `len`

## Project Structure

- `lexer/`: Lexical analyzer for tokenizing source code.
- `parser/`: Recursive descent + Pratt parser for AST generation.
- `ast/`: Abstract Syntax Tree node definitions.
- `object/`: Runtime object system and environment handling.
- `evaluator/`: Tree-walking interpreter and built-in functions.
- `repl/`: Interactive shell implementation.
- `token/`: Token definitions and constants.
- `samples/`: Example programs.

## Future Plans

- Ownership and Borrow Checker
- Modules and Imports
- Concurrency (Goroutines/Channels)
- Structs and Traits
- Bytecode VM or Native Compilation
