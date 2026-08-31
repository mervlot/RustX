# RustX Starter

![RustX Logo](assets/rustx-logo.png)

RustX Starter is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and future extensibility, RustX Starter currently operates as a lightweight scripting and systems hybrid language while maintaining an extensible architecture for future advanced features.

---

## Language Highlights

- **Syntax Blend:** Combines Go's concise `func`, `var`, `:=`, and loop constructs with Rust-like modern aesthetics.
- **Optional Semicolons & Typing:** Flexible statement termination and optional type annotations in function signatures and variable declarations.
- **Top-Level Entry Point:** Automatically runs `main()` if defined at the top-level script scope.
- **Standard Library:** Built-in support for standard output (`print`, `println`), math functions, string manipulation, time utilities, and filesystem operations.
- **Human-Readable Diagnostics:** Clear location-based `SyntaxError` and `RuntimeError` messages with line and column precision.

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.20 or later.

### Building the Compiler/Interpreter

From the repository root, run:

```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

### 1. Execute Scripts Directly
```bash
./rustx run samples/hello.rx
```
or using shorthand:
```bash
./rustx samples/hello.rx
```

### 2. Interactive REPL
Launch the REPL prompt:
```bash
./rustx repl
```
or run `./rustx` without arguments:
```bash
$ ./rustx
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> println(x * 2)
20
```

### 3. Syntax Validation (Build Check)
Check syntax without executing the file:
```bash
./rustx build samples/factorial.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## Language Syntax Guide

### Variables
```go
var age int = 16   // Explicit variable declaration
var city = "Lagos" // Type inferred variable
name := "RustX"    // Shorthand assignment
```

### Control Flow
```go
// If-Else
if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

// For Loop (3-part)
for i := 0; i < 5; i++ {
    println("Count:", i)
}

// While Loop
count := 3
while count > 0 {
    println(count)
    count--
}
```

### Functions & Recursion
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5:", factorial(5))
}
```

---

## Architecture Overview

RustX Starter is designed using a clean, modular architecture:

```
├── ast/        # Abstract Syntax Tree node definitions
├── errors/     # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/  # AST evaluation engine and standard library built-ins
├── lexer/      # Lexical analyzer / tokenizer with position tracking
├── object/     # Runtime object representation and Environment scoping
├── parser/     # Recursive descent Pratt parser
├── repl/       # Interactive REPL with multiline support
├── samples/    # Example .rx programs
├── token/      # Token types and definitions
└── main.go     # CLI entry point
```

### Future Extensibility
Though currently functioning as an interpreted scripting language, AST nodes, token structures, and environment scopes were designed so future versions can introduce:
- Ownership & Borrow Checking
- Native Compilation / Bytecode VM
- Structs, Enums, and Interfaces / Traits
- Concurrency (`go` / `async`/`await`)
- Module / Import system

---

## License

MIT
