# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntax and strengths of both languages to achieve simplicity, speed, readability, and future extensibility.

> **Note on Architecture:** RustX Starter serves as the foundational scripting + systems hybrid language. Memory safety systems (such as borrow checking, ownership models, and lifetimes) are intentionally omitted in this foundation phase, but the modular lexer/parser/AST/runtime architecture is engineered so those features can be added in future iterations without refactoring the compiler/interpreter core.

---

## Features

- **Blended Syntax:** Combines Go's `func`, variable declaration `var`, and shorthand `:=` with Rust-like modern expressiveness.
- **Optional Semicolons:** Semicolons are supported but optional at statement ends.
- **Control Flow:** `if` / `else`, `while`, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Data Types:** `int`, `float`, `bool`, `string`, `array`, and functions.
- **Built-in Standard Library:** Includes core functions for `strings` (len, contains, split, join, to_upper, to_lower), `math` (abs, pow, sin, cos, tan), `time` (now, sleep), `filesystem` (read_file, write_file), and printing (`println`, `print`).
- **Automatic Entrypoint:** Automatically invokes `main()` if defined at script top-level.
- **Multiline REPL:** Interactive shell supporting multiline blocks and persistent environment state across commands.
- **Human-Readable Diagnostics:** Clear position-aware `SyntaxError` and `RuntimeError` messages.

---

## Installation & Build Instructions

### Prerequisites

- [Go](https://golang.org/) 1.18 or higher.

### Building RustX

Build the executable binary from the project root:

```bash
go build -o rustx main.go
```

Verify installation:

```bash
./rustx version
```

---

## CLI Usage

RustX Starter provides a CLI tool (`rustx`):

```bash
# Start the multiline interactive REPL
rustx repl
# or simply
rustx

# Execute a RustX source script
rustx run samples/hello.rx

# Direct shorthand execution
rustx samples/hello.rx

# Syntax validation / Build check without executing
rustx build samples/hello.rx

# View version
rustx version
```

---

## Language Syntax Guide

### Variables and Assignment

```rustx
// Mutable variable with type annotation
var age int = 16

// Mutable variable with type inference
var count = 0

// Shorthand variable assignment
name := "Muhammad"

// Reassignment
count = count + 1
```

### Functions

Functions use the `func` keyword, support parameter and return type annotations, and can be recursive:

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

func main() {
    println("Factorial of 5:", factorial(5))
}
```

### Control Flow

Unparenthesized conditions are fully supported:

```rustx
if name == "RustX" {
    println("Language initialized")
} else {
    println("Other language")
}

for i := 0; i < 5; i++ {
    println("Count:", i)
}

while count > 0 {
    count--
}
```

### Standard Library Functions

```rustx
// Printing
println("Hello", "World")
print("No newline")

// Strings
len("RustX")               // 5
contains("RustX", "ust")  // true
to_upper("rustx")          // "RUSTX"
to_lower("RUSTX")          // "rustx"
split("a,b,c", ",")        // ["a", "b", "c"]
join(["a", "b"], "-")      // "a-b"

// Math
abs(-42)                   // 42
pow(2, 3)                  // 8.0
sin(0)                     // 0.0

// Filesystem
write_file("data.txt", "RustX output")
content := read_file("data.txt")
```

---

## Architecture & Folder Structure

```
.
├── ast/         # Abstract Syntax Tree nodes implementing GetToken() and String()
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime evaluator and built-in standard library functions
├── lexer/       # Lexical analyzer / tokenizer with comment & string escape support
├── object/      # Runtime object system and Environment scope management
├── parser/      # Recursive-descent Pratt parser with operator precedence
├── repl/        # Multiline REPL prompt with persistent environment
├── samples/     # Sample .rx scripts demonstrating language capabilities
├── token/       # Token types and definitions
├── main.go      # CLI entrypoint
├── go.mod       # Go module declaration
└── README.md    # Documentation
```

---

## Running Unit Tests

Run the full test suite across all packages:

```bash
go test ./...
```
