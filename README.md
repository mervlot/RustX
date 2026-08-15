# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and future extensibility, RustX Starter blends the expressive ergonomics of Rust with the straightforward simplicity of Go.

> **Note on Architecture & Future Safety System:**
> RustX Starter deliberately does **NOT** implement a borrow checker, ownership model, or lifetimes in this foundation phase. Its modular architecture (Lexer, Parser, AST, Evaluator, Runtime) is strictly designed so memory safety and systems features can be integrated in future phases without rewriting the interpreter/compiler.

---

## Language Features

- **Blended Syntax:** `func` for functions, `var` for mutable variables, `:=` shorthand declaration, curly brace blocks, optional semicolons.
- **Primitive Types:** `int`, `float`, `bool`, `string`, plus basic arrays `[...]`.
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), and post increment/decrement (`++`, `--`).
- **Control Flow:** `if` / `else if` / `else`, `while` loops, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`).
- **Functions:** Parameters, optional type annotations, recursion, first-class/anonymous functions.
- **Built-in Standard Library:**
  - **Output:** `print`, `println`
  - **Strings:** `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math:** `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time:** `sleep`, `now`
  - **Filesystem:** `read_file`, `write_file`
- **Position-Aware Diagnostics:** Friendly `SyntaxError` and `RuntimeError` messages containing exact line and column numbers.

---

## Directory Structure

```text
rustx/
├── ast/         # AST node definitions implementing Node, Statement, Expression
├── errors/      # Position-aware SyntaxError and RuntimeError types
├── evaluator/   # Tree-walking interpreter and built-in standard library functions
├── lexer/       # Tokenizer supporting single/multi-line comments and escape characters
├── object/      # Runtime object system and lexical environment scoping
├── parser/      # Recursive Descent + Pratt parser with position tracking
├── repl/        # Interactive REPL prompt supporting multiline blocks
├── samples/     # Example .rx scripts
├── token/       # Token definitions and keyword map
├── main.go      # CLI entry point
├── README.md    # Documentation
└── go.mod       # Go module declaration
```

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.20+ installed.

### Building RustX

Build the `rustx` executable from the root directory:

```bash
go build -o rustx main.go
```

---

## Usage & CLI

### 1. Execute a Script
Run `.rx` files using `run` or direct execution shorthand:

```bash
./rustx run samples/hello.rx
# or
./rustx samples/hello.rx
```

### 2. Interactive REPL
Start the interactive shell:

```bash
./rustx repl
# or simply
./rustx
```

REPL prompt example:
```text
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+C to quit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> for i := 0; i < 3; i++ {
...     println(i)
... }
0
1
2
```

### 3. Syntax Validation
Validate syntax without executing the program:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Info
Display language version:

```bash
./rustx version
```

---

## Syntax Examples

### Variable Declarations

```rustx
var age int = 16
var city = "New York"
name := "Muhammad"
```

### Functions & Recursion

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

println(factorial(5)) // Outputs: 120
```

### Loops and Control Flow

```rustx
for i := 1; i <= 5; i++ {
    println("Iteration:", i)
}

var x = 0
while x < 3 {
    println("x is", x)
    x++
}
```

---

## Running Unit Tests

Run all unit tests across all packages:

```bash
go test ./...
```
