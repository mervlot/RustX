<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="220" />
</p>

# RustX Starter

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**, focused on simplicity, speed, readability, and future extensibility. It bridges the gap between lightweight scripting and systems programming.

> **Design Directive Note**: Advanced memory safety systems (ownership, lifetimes, borrow checker) are intentionally omitted in this initial foundation phase. The AST and runtime architecture are designed modularly so these systems can be added seamlessly in future versions without rewriting the compiler or interpreter.

---

## Features

- **Blended Syntax**: Familiar and concise syntax blending Go (`func`, `var`, `:=`, optional semicolons) and Rust expressions.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Relational (`==`, `!=`, `<`, `<=`, `>`, `>=`), Logical short-circuiting (`&&`, `||`), Increment/Decrement (`++`, `--`).
- **Functions**: First-class functions, closures, optional type annotations, recursion, and auto-executed `main()` entry points.
- **Control Flow**: `if` / `else if` / `else`, `while`, and 3-part or conditional `for` loops.
- **Standard Library**: Built-in modules for printing (`println`, `print`), collections (`len`, `push`), strings (`str_to_upper`, `str_to_lower`, `str_contains`, `str_split`), math (`math_sqrt`, `math_abs`, `math_pow`), time (`time_now_unix`, `time_sleep_ms`), and filesystem (`read_file`, `write_file`).
- **Interactive REPL**: Multiline input support, variable/function preservation across prompts, human-readable error diagnostics.

---

## Project Architecture

```
rustx/
├── ast/         # Abstract Syntax Tree definitions for statements and expressions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walking interpreter and built-in standard library
├── lexer/       # Lexer / tokenizer with position tracking and comment handling
├── object/      # Runtime object model and environment lexical scoping
├── parser/      # Recursive descent & Pratt parser with operator precedences
├── repl/        # Interactive multiline read-eval-print-loop
├── samples/     # Example programs (.rx files)
├── token/       # Token types and reserved keyword tables
├── main.go      # CLI entrypoint (run, repl, version, build)
└── go.mod       # Go module configuration
```

---

## Installation & Building

RustX requires **Go 1.18+** (tested on Go 1.24+).

### Build from source:

```bash
go build -o rustx main.go
```

---

## CLI Usage

### 1. Run a RustX Program
Execute `.rx` script files directly:
```bash
./rustx run samples/hello.rx
# or shorthand:
./rustx samples/hello.rx
```

### 2. Interactive REPL
Launch the interactive prompt:
```bash
./rustx repl
# or simply:
./rustx
```

Example session:
```text
RustX Starter v0.1
Type statements or expressions. Press Ctrl+C or Ctrl+D to exit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func double(n int) int {
...     return n * 2
... }
>>> double(21)
42
```

### 3. Syntax Validation / Build
Validate file syntax without executing:
```bash
./rustx build samples/factorial.rx
```

### 4. Version
Display version info:
```bash
./rustx version
```

---

## Syntax & Examples

### Hello World
```go
func main() {
    println("Hello World from RustX Starter!")
}
```

### Recursion & Functions
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5)) // 120
```

### Variables and Control Flow
```go
var count = 0
name := "RustX"

for i := 0; i < 5; i++ {
    count++
}

if count == 5 && name == "RustX" {
    println("Initialized successfully")
}
```

### Builtin Standard Library
```go
// Strings
println(str_to_upper("rustx"))

// Math
println(math_sqrt(64)) // 8

// Filesystem
write_file("greeting.txt", "Hello from RustX!")
content := read_file("greeting.txt")
print(content)
```

---

## Running Tests

Run the test suite across all modules:
```bash
go test ./...
```
