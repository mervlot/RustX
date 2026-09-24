# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntactic elegance of Rust and Go with a fast interpreter and standard library, designed as an extensible foundation for future systems language features (such as ownership, borrow checking, and static compilation) without requiring a core compiler rewrite.

---

## Key Language Features

- **Blended Go & Rust Syntax**: `func` for functions, `var` for declarations, `:=` for shorthand variable initialization, optional type annotations, and optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `>`, `<`, `>=`, `<=`), logic (`&&`, `||`), prefix (`!`, `-`), postfix (`++`, `--`).
- **Control Flow**: `if` / `else if` / `else`, `while` loops, Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Functions & Recursion**: First-class function literals, named functions, lexical environment scoping, recursive function calls.
- **Built-in Standard Library**:
  - **I/O & Printing**: `println()`, `print()`.
  - **Strings**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`.
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`.
  - **Time**: `now()`, `sleep()`.
  - **Filesystem**: `read_file()`, `write_file()`.
- **Diagnostics**: Position-aware human-readable error messages with exact line and column locations (`SyntaxError` and `RuntimeError`).

---

## Directory Structure & Architecture

```text
rustx/
├── token/         Token definitions and keyword lookup
├── lexer/         Lexical scanner with line & column tracking, comments & escapes
├── ast/           AST nodes for statements and Pratt expressions
├── parser/        Recursive descent & Pratt expression parser
├── object/        Runtime object types and Environment scope hierarchy
├── evaluator/     Tree-walking interpreter & standard library built-ins
├── repl/          Interactive prompt with multiline input tracking
├── errors/        Standardized SyntaxError and RuntimeError diagnostics
├── samples/       Demonstration scripts (.rx)
├── main.go        CLI entry point
└── README.md
```

---

## Installation & Build Instructions

### Prerequisites
- Go 1.18 or higher

### Building RustX Starter
```bash
go build -o rustx main.go
```

### Running Tests
To run all unit tests across the project:
```bash
go test ./...
```

---

## CLI Usage

```bash
# Run a script directly
./rustx run samples/hello.rx

# Shorthand execution
./rustx samples/factorial.rx

# Start interactive REPL prompt
./rustx repl

# Validate file syntax without executing
./rustx build samples/std_demo.rx

# Display version information
./rustx version
```

---

## Language Syntax Guide

### Variables
```rustx
var count = 10
name := "Muhammad"
```

### Functions
```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("5! =", factorial(5))
```

### Control Flow
```rustx
// If - Else
if x > 10 {
    println("Large")
} else {
    println("Small")
}

// For Loop (3-part)
for i := 0; i < 5; i++ {
    println(i)
}

// While Loop
while x > 0 {
    x--
}
```

### Standard Library Examples
```rustx
// Strings
str := "hello rustx"
println(to_upper(str)) // "HELLO RUSTX"

// Filesystem
write_file("data.txt", "RustX Starter Storage Test")
content := read_file("data.txt")
println("File content:", content)
```

---

## Future Extensibility

RustX Starter is designed with modular AST nodes, position tracking on every node, and clean object evaluation. This architecture allows future releases to introduce:
- Ownership & Borrow Checker
- Memory Lifetimes
- Structs, Enums, and Pattern Matching
- Interfaces / Traits
- Concurrency & Async/Await
- Bytecode VM / LLVM native compilation
