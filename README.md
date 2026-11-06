# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="220" />
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**, focused on simplicity, speed, readability, and future extensibility. It blends the ergonomics and familiarity of Go and Rust into a lightweight scripting + systems hybrid language.

> **Note on Safety Model:** Advanced memory safety systems (such as a borrow checker, ownership semantics, and lifetimes) are intentionally deferred in this foundation phase. The AST and runtime architecture are designed with clean boundaries so those systems, as well as a bytecode VM or native compiler backend, can be added in future iterations without rewriting the parser or frontend.

---

## Features

- **Blended Go & Rust Syntax**: `func` for functions, `var` and `:=` shorthand, curly brace blocks, optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`, `!`), Postfix increment/decrement (`++`, `--`).
- **Control Flow**: `if` / `else`, `while`, and 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`).
- **Functions & Recursion**: First-class user-defined functions with support for recursion, closures, and lexical scoping. If a parameterless `main()` function is defined in a script, it is automatically called as the entrypoint.
- **Built-in Standard Library**:
  - **Console I/O**: `print(...)`, `println(...)`
  - **Collection Utilities**: `len(obj)`, `push(arr, item)`
  - **Strings**: `str_contains(s, substr)`, `str_split(s, sep)`, `str_join(arr, sep)`, `str_upper(s)`, `str_lower(s)`, `str_trim(s)`
  - **Math**: `abs(x)`, `sqrt(x)`, `pow(base, exp)`, `min(a, b)`, `max(a, b)`, `round(x)`, `floor(x)`, `ceil(x)`
  - **Time**: `time_now()`, `sleep_ms(ms)`
  - **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Human-Readable Diagnostics**: Position-aware error reporting with line and column numbers.
- **Interactive Multiline REPL**: Interactive environment preserving variables and functions across lines, supporting nested blocks and braces.

---

## Architecture

The project adheres to a clean, modular architecture:

```
├── ast/         # Abstract Syntax Tree definitions and node interfaces
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostic definitions
├── evaluator/   # Tree-walking interpreter and built-in standard library functions
├── lexer/       # Lexical scanner converting source code into tokens with line/col tracking
├── object/      # Runtime object representations and lexical environment scope chains
├── parser/      # Recursive descent / Pratt expression parser
├── repl/        # Interactive multiline REPL implementation
├── samples/     # Example programs demonstrating language capabilities
├── token/       # Token types and keyword mappings
└── main.go      # CLI entrypoint and argument dispatcher
```

---

## Installation & Building

RustX Starter is written in standard Go (Go 1.20+ recommended).

### Build Binary

```bash
go build -o rustx main.go
```

### Run Tests

```bash
go test -v ./...
```

---

## Usage

### 1. Interactive REPL

Start the REPL prompt:

```bash
./rustx repl
# or simply
./rustx
```

Example session:

```
RustX Starter v0.1
Type RustX expressions or statements. Press Ctrl+C or Ctrl+D to exit.
>>> name := "RustX"
>>> println("Hello from", name)
Hello from RustX
>>> func fib(n int) int {
...     if n <= 1 {
...         return n
...     }
...     return fib(n - 1) + fib(n - 2)
... }
>>> fib(7)
13
```

### 2. Running Scripts

Execute `.rx` files:

```bash
./rustx run samples/hello.rx
# or shorthand:
./rustx samples/hello.rx
```

### 3. Syntax Validation (`build`)

Check a script for syntax errors without executing:

```bash
./rustx build samples/factorial.rx
```

### 4. Check Version

```bash
./rustx version
```

---

## Code Examples

### Functions & Recursion

```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### Loops

```go
for i := 0; i < 5; i++ {
    println("Iteration:", i)
}

count := 3
while count > 0 {
    println("Countdown:", count)
    count--
}
```

### Arrays and Standard Library

```go
arr := [10, 20, 30]
push(arr, 40)
println("Array length:", len(arr))

msg := "  RustX Starter  "
println(str_upper(str_trim(msg)))
```
