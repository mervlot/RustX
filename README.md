# RustX Starter

RustX Starter is a lightweight, compiled/interpreted scripting and systems hybrid programming language inspired by **Rust** and **Go**. It is designed for simplicity, readability, speed, and future extensibility.

---

## 🌟 Architecture Overview

RustX Starter is implemented in clean, modular Go with an extensible package architecture:

```text
rustx/
├── token/       # Token types and position tracking (Line, Column)
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── lexer/       # Lexical scanner for tokens, comments, and strings
├── ast/         # Abstract Syntax Tree nodes implementing Node interface
├── parser/      # Recursive descent / Pratt parser
├── object/      # Runtime objects and environment with lexical scope
├── evaluator/   # AST evaluator and standard library built-ins
├── repl/        # Interactive REPL with multiline support
├── samples/     # Example programs written in .rx
└── main.go      # CLI entry point
```

---

## 🚀 Getting Started & Building

### Prerequisites
- Go 1.20+ installed

### Build the Executable
```bash
go build -o rustx main.go
```

### Run All Unit Tests
```bash
go test ./...
```

---

## 💻 CLI Usage

### Execute Scripts
```bash
# Run with 'run' command
./rustx run samples/hello.rx

# Shorthand execution
./rustx samples/factorial.rx
```

### Interactive REPL
```bash
./rustx repl
# or simply
./rustx
```

Example REPL session:
```text
RustX Starter v0.1
Type RustX expressions or commands. Press Ctrl+C to exit.
>>> x := 10
>>> println(x * 2)
20
```

### Syntax Validation (Build Check)
```bash
./rustx build samples/factorial.rx
```

### Display Version
```bash
./rustx version
```

---

## 📝 Syntax & Language Features

### Variables & Mutability
```rustx
var age int = 25
var city = "San Francisco"
name := "RustX"
```

### Data Types
- `int` (64-bit integers)
- `float` (64-bit floating point)
- `bool` (`true`, `false`)
- `string` (UTF-8 string literals with escape sequences)
- `[T]` (arrays)

### Control Flow
```rustx
if name == "RustX" {
    println("Welcome to RustX!")
} else {
    println("Other language")
}

for i := 0; i < 5; i++ {
    println(i)
}

count := 0
while count < 3 {
    count++
}
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
```

### Built-in Standard Library
- **IO**: `print()`, `println()`
- **Strings**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
- **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
- **Time**: `now()`, `sleep()`
- **Filesystem**: `read_file()`, `write_file()`
