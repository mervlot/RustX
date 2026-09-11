# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a clean, modern, lightweight compiled/interpreted programming language that blends the elegant syntax of **Rust** and **Go**. It is designed with simplicity, speed, readability, and future extensibility in mind.

> **Note on Architecture:** RustX Starter serves as the solid interpreter foundation. In this phase, advanced memory safety systems (ownership, borrow checker, lifetimes) are intentionally deferred, but the modular AST and environment architecture are built to support future addition of bytecode VM, native compilation, generics, structs, enums, traits, and an ownership borrow checker without a full rewrite.

---

## Language Highlights

- **Syntax Blend:** `func` for function declarations, `var` and `:=` shorthand for variable bindings, strong typing hints, optional semicolons, curly brace control blocks.
- **Data Types:** `int` (64-bit), `float` (64-bit), `bool`, `string`, `Array`.
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||` with short-circuit evaluation), Postfix (`++`, `--`).
- **Control Flow:** `if` / `else if` / `else`, `while` loops, and 3-part / single-condition `for` loops.
- **Functions:** Named and anonymous functions, optional parameter and return type annotations, recursive support.
- **Standard Library:** Built-in support for printing (`print`, `println`), strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem operations (`read_file`, `write_file`).
- **Position-Aware Diagnostic Errors:** Clean, human-readable error messages detailing line and column location without interpreter panics.

---

## Directory & Modular Architecture

The repository is organized cleanly into modular packages:

```
├── token/       # Token types, constants, keyword map, and source position metadata
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── lexer/       # Lexical analyzer with comment skipping and string escape decoding
├── ast/         # Abstract Syntax Tree nodes implementing Node, Statement, Expression, GetToken()
├── parser/      # Pratt precedence-based parser & recursive descent control flow parser
├── object/      # Runtime object types and Environment with scope nesting
├── evaluator/   # Evaluator engine, short-circuiting, and standard library builtins
├── repl/        # Interactive REPL prompt with multiline input & balanced block detection
├── samples/     # Sample RustX scripts (.rx)
├── main.go      # CLI entry point (run, build, repl, version)
└── README.md
```

---

## Getting Started & Build Instructions

### Prerequisites
- [Go](https://golang.org/) 1.18 or higher installed on your machine.

### Building the CLI
To compile the `rustx` executable binary:
```bash
go build -o rustx main.go
```

---

## CLI & Usage Guide

RustX Starter supports multiple CLI commands and direct file execution shorthand:

### 1. Running a Script
```bash
./rustx run samples/hello.rx
```
or shorthand:
```bash
./rustx samples/hello.rx
```

### 2. Interactive REPL Prompt
Launch the multiline interactive prompt:
```bash
./rustx repl
```
or simply run without arguments:
```bash
./rustx
```

Example REPL session:
```text
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> func double(n int) int {
...     return n * 2
... }
>>> double(x)
20
```

### 3. Syntax Validation (`build`)
To check syntax errors in a `.rx` file without executing runtime logic:
```bash
./rustx build samples/factorial.rx
```

### 4. Version Check
```bash
./rustx version
```

---

## Code Examples

### Example 1: Hello World
`samples/hello.rx`
```rustx
func main() {
    println("RustX Starter")
}
```

### Example 2: Recursive Factorial
`samples/factorial.rx`
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

### Example 3: Control Flow & Standard Library
`samples/control_flow.rx`
```rustx
func main() {
    name := "RustX"

    if name == "RustX" {
        println("Language initialized successfully!")
    }

    // Loops
    var total int = 0
    for i := 1; i <= 10; i++ {
        total = total + i
    }
    println("Sum of 1..10 is:", total)

    // Standard Library
    text := "hello world"
    println("Upper:", to_upper(text))
    println("Contains 'world':", contains(text, "world"))
    println("2^10 =", pow(2, 10))
}
```

---

## Testing

To run the complete unit test suite across all packages:
```bash
go test -v ./...
```
