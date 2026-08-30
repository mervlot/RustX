# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180">
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

The language operates as a lightweight scripting + systems hybrid language.

---

## Key Features

- **Blended Syntax**: Integrates `func` for function declarations, `var` for mutable variables, and `:=` for shorthand variable declarations.
- **Flexible Execution**: Run scripts directly via CLI or interactively through a multiline-capable REPL.
- **Control Flow**: Supports `if`, `else`, `while`, and Go-style `for` loops with optional parentheses.
- **Data Types & Literals**: Core primitives (`int`, `float`, `bool`, `string`), arrays (`[1, 2, 3]`), and string escape sequences.
- **Standard Library**: Built-in utilities for string handling (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem ops (`read_file`, `write_file`).
- **Human-Readable Diagnostics**: Precise `SyntaxError` and `RuntimeError` reporting with line and column tracking.
- **Extensible Architecture**: Modular lexer, Pratt parser, AST, and evaluator designed for future integration of ownership models, borrow checkers, and compiler targets.

---

## Directory Architecture

```
rustx/
├── token/       # Token types, keywords, and position definitions
├── lexer/       # Lexical analyzer with line/col tracking and comment skipping
├── ast/         # Extensible Abstract Syntax Tree (AST) node representations
├── parser/      # Recursive descent / Pratt parser with precedence handling
├── object/      # Runtime object model and environment with lexical scoping
├── evaluator/   # AST evaluator and built-in standard library functions
├── errors/      # Position-aware SyntaxError and RuntimeError formatters
├── repl/        # Stateful, multiline-capable interactive shell
├── samples/     # Example .rx script files
├── main.go      # CLI entry point
├── go.mod       # Module definition
└── README.md
```

---

## Quick Start & Installation

### Build Binary
```bash
go build -o rustx main.go
```

### Run Tests
```bash
go test ./...
```

---

## CLI Usage

### Execute Script File
```bash
./rustx run samples/hello.rx
# Or using shorthand
./rustx samples/hello.rx
```

### Interactive REPL
```bash
./rustx repl
# Or run without arguments
./rustx
```

### Check Syntax (Build Check)
```bash
./rustx build samples/factorial.rx
```

### Display Version
```bash
./rustx version
```

---

## Language Syntax Guide

### Functions & Main Entrypoint
```rustx
func add(a int, b int) int {
    return a + b
}

func main() {
    println("Hello World")
    println(add(5, 10))
}
```

### Variable Declarations
```rustx
var count = 10
name := "RustX"
```

### Recursion & Control Flow
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### Loops
```rustx
for i := 0; i < 5; i++ {
    println(i)
}
```

---

## Future Compatibility & Roadmap

RustX Starter is designed as a foundational stepping stone. The core architecture isolates execution from syntax definition, making it ready to support:
- Ownership & Borrow Checker rules
- Structs, Enums, & Pattern Matching
- Generics & Traits/Interfaces
- Module & Package Manager system
- Native Compilation / Bytecode VM backend
