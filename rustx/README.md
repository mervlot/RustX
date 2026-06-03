# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go. It is designed to be a lightweight hybrid of a scripting and systems language, with an architecture that allows for future features like a borrow checker, ownership model, and native compilation.

## Features

- **Variables:** Supports `var` declarations and `:=` shorthand.
- **Data Types:** `int`, `float`, `bool`, and `string`.
- **Operators:** Standard arithmetic, comparison, and logical (`&&`, `||`) with short-circuiting.
- **Functions:** Parameters, return types, and recursion support.
- **Control Flow:** `if`, `else`, `while`, and `for` loops.
- **Built-in Functions:** Math, strings, time, and filesystem basics.
- **REPL:** Interactive prompt with multiline support.

## Project Structure

- `/lexer`: Tokenizes the source code.
- `/token`: Defines the language tokens.
- `/parser`: Recursive descent parser (Pratt parsing for expressions).
- `/ast`: Abstract Syntax Tree nodes.
- `/object`: Object system and runtime environment.
- `/evaluator`: Tree-walking interpreter.
- `/repl`: Interactive REPL implementation.
- `main.go`: CLI entry point.

## Installation & Usage

### Prerequisites
- [Go](https://golang.org/doc/install) (1.16 or later)

### Build
To build the `rustx` executable:
```bash
cd rustx
go build -o rustx main.go
```

### Run a Script
```bash
./rustx run examples/demo.rx
```

### Start REPL
```bash
./rustx repl
```

### Version
```bash
./rustx version
```

## Syntax Example

```rust
func factorial(n) {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:")
println(factorial(5))

x := 10
x = 20 // Re-assignment

for i := 0; i < 5; i++ {
    println(i)
}
```
