# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

A minimal, clean, scalable interpreted/compiled hybrid programming language inspired by Rust and Go syntax. Built with simplicity, speed, readability, and future extensibility in mind.

---

## Overview

**RustX Starter** combines the syntactic elegance of Rust and Go:
- **Rust-inspired** function declarations (`func`), strong expression semantics, and future compatibility for ownership models.
- **Go-inspired** simplicity, variable definitions (`var` and shorthand `:=`), control flow structures, and clean modular architecture.

> **Note on Architecture:** RustX Starter currently operates as a lightweight scripting/systems hybrid language. It does **not** include a borrow checker, lifetimes, or ownership enforcement yet—however, the AST and parser architecture are designed so these systems can be added cleanly in future iterations.

---

## Features

- **Variables & Mutability:** `var x = 10` or shorthand `x := 10`
- **Data Types:** `int`, `float`, `bool`, `string`, `array`
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`)
- **Functions:** Named and anonymous function declarations, parameters, return type annotations, recursive support
- **Control Flow:** `if` / `else if` / `else`, `while` loops, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`)
- **Standard Library:** Built-in support for string operations, math functions, time, and basic filesystem I/O
- **REPL & Script Execution:** Interactive prompt with multiline input handling and script execution from CLI
- **Human-Readable Diagnostics:** Structured position-aware `SyntaxError` and `RuntimeError` reporting

---

## Directory Structure

```text
rustx/
├── ast/         # AST node definitions (statements, expressions, types)
├── errors/      # Position-aware error reporting (SyntaxError, RuntimeError)
├── evaluator/   # Interpreter, tree-walking evaluation engine & builtins
├── lexer/       # Lexical analysis and tokenizer
├── object/      # Runtime object system and lexical environment scoping
├── parser/      # Recursive descent + Pratt parser engine
├── repl/        # Interactive REPL with multiline input support
├── samples/     # Example .rx scripts
├── token/       # Token types and definitions
├── go.mod       # Go module declaration
├── main.go      # CLI entry point
└── README.md    # Documentation
```

---

## Getting Started

### Prerequisites
- Go 1.20+ (1.24+ recommended)

### Build
Compile the `rustx` executable:
```bash
go build -o rustx main.go
```

### Run Tests
Execute unit test suite across all packages:
```bash
go test ./...
```

---

## CLI Usage

```bash
# Run a script file
./rustx run samples/hello.rx

# Shorthand execution
./rustx samples/factorial.rx

# Interactive REPL
./rustx repl

# Check / validate syntax of a file
./rustx build samples/stdlib_demo.rx

# View version information
./rustx version
```

---

## Code Examples

### 1. Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Recursion & Factorial (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5 is:", factorial(5))
}
```

### 3. Standard Library Demo (`samples/stdlib_demo.rx`)
```rustx
func main() {
    name := "RustX Starter"

    if contains(name, "RustX") {
        println("Found RustX in name!")
    }

    upperName := to_upper(name)
    println("Uppercase:", upperName)

    parts := split(name, " ")
    println("Joined again:", join(parts, " - "))

    println("Absolute value of -42:", abs(-42))
    println("2 raised to power 8:", pow(2, 8))
}
```

---

## Future Roadmap

RustX Starter is designed as a foundational stepping stone towards a complete compiler and runtime ecosystem:
- [ ] Ownership & Borrow Checker
- [ ] Module system (`import`)
- [ ] Structs, Enums, and Traits / Interfaces
- [ ] Generics
- [ ] Concurrency (`async` / `await` or goroutine-style tasks)
- [ ] Bytecode VM or Native LLVM / Cranelift Compiler Backend
