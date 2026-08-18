# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and extensible programming language inspired by **Rust** and **Go**. It focuses on simplicity, readability, execution speed, and future extensibility.

---

## Key Features

- **Blended Syntax**: Combines Go's ergonomic shorthand declarations (`:=`, `var`, `func`) with Rust-inspired syntax and optional type annotations.
- **Dual Execution**: Run `.rx` scripts directly or launch an interactive REPL with multiline support.
- **Strong Foundation**: Built with a clean, modular architecture (Lexer, Pratt Parser, AST, Evaluator runtime, and Environment scoping).
- **Extensible Architecture**: Designed specifically so future features like ownership/borrow-checking systems, bytecode compilation, structs, enums, generics, and concurrency can be added without architectural rewrites.
- **Human-Readable Diagnostics**: Clean position-aware syntax and runtime error reports.
- **Standard Library**: Built-in utilities for printing, math operations, string manipulation, time, and filesystem operations.

---

## Folder Structure

```
.
├── ast/         # Abstract Syntax Tree node definitions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walking interpreter & standard library built-ins
├── lexer/       # Lexical analyzer / tokenizer with position tracking
├── object/      # Runtime object system and enclosed environment scopes
├── parser/      # Recursive descent + Pratt expression parser
├── repl/        # Interactive REPL prompt with multiline support
├── samples/     # Example .rx scripts
├── token/       # Token types, constants, and keyword lookup
├── main.go      # CLI entry point
└── README.md    # Documentation
```

---

## Installation & Building

### Prerequisites
- [Go 1.20+](https://golang.org/dl/)

### Build Binary
To build the standalone executable:

```bash
go build -o rustx main.go
```

To run unit tests across all packages:

```bash
go test ./...
```

---

## CLI Usage

RustX Starter provides simple and flexible CLI commands:

```bash
# Execute a script
./rustx run samples/hello.rx

# Direct script execution shorthand
./rustx samples/factorial.rx

# Start the interactive REPL
./rustx repl
# OR simply:
./rustx

# Validate file syntax without executing
./rustx build samples/control_flow.rx

# Print language version
./rustx version
```

---

## Language Syntax Guide

### Variables & Constants
Variables are mutable by default. Support for `var` and `:=` shorthand assignment:

```rustx
var age = 16
name := "Muhammad"
```

### Data Types
- `int`: `42`
- `float`: `3.14159`
- `bool`: `true` / `false`
- `string`: `"Hello RustX\n"`
- `array`: `[1, 2, 3]`

### Functions
Functions use the `func` keyword and support optional parameter and return type annotations:

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

println("Factorial of 5:", factorial(5))
```

### Control Flow
Conditionals and loops do not require parentheses around headers:

```rustx
// If - Else
if name == "Muhammad" {
    println("Welcome!")
} else {
    println("Hello Guest")
}

// Three-part For loop
for i := 0; i < 5; i++ {
    println(i)
}

// While loop
count := 3
while count > 0 {
    println("Count:", count)
    count--
}
```

### Standard Library Built-ins
- **I/O**: `print(...)`, `println(...)`
- **Strings**: `len(s)`, `contains(s, substr)`, `split(s, sep)`, `join(array, sep)`, `to_upper(s)`, `to_lower(s)`
- **Math**: `abs(x)`, `pow(base, exp)`, `sin(x)`, `cos(x)`, `tan(x)`
- **Time**: `now()`, `sleep(ms)`
- **Filesystem**: `read_file(path)`, `write_file(path, content)`

---

## Future Compatibility & Roadmap

RustX Starter is designed as a foundational stepping stone. The architecture isolates lexical analysis, AST construction, and environment resolution, ensuring seamless future enhancements:

1. **Ownership & Borrow Checker**: Linear type checks during semantic AST pass.
2. **Structs, Enums & Pattern Matching**: New AST nodes and runtime object layouts.
3. **Bytecode VM / Native Compiler**: Replace tree-walker evaluator with bytecode compiler and virtual machine.
4. **Concurrency**: Goroutine-style lightweight tasks or async/await syntax.
