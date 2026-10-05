# RustX Starter

![RustX Logo](assets/rustx-logo.png)

RustX Starter is a lightweight, clean, and scalable scripting + systems hybrid programming language inspired by **Rust** and **Go**. Designed with simplicity, speed, readability, and future extensibility in mind, RustX Starter provides a modern language foundation without immediate memory management overheads (ownership/borrow checker are deferred to future phases).

---

## 🌟 Language Overview

- **Extension:** `.rx`
- **Keywords:** `func`, `var`, `if`, `else`, `while`, `for`, `return`, `true`, `false`
- **Shorthand Assignment:** `:=`
- **Types:** `int`, `float`, `bool`, `string`, `array` (optional type annotations supported)
- **Syntax Highlights:** Curly brace blocks, optional semicolons, Go/Rust hybrid control flow
- **Built-in Standard Library:**
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time & FS: `sleep`, `now`, `read_file`, `write_file`

---

## 🏗️ Project Architecture

```text
rustx/
├── token/       # Token types and position tracking (line/column)
├── lexer/       # Lexical scanner for keywords, identifiers, numbers, strings & comments
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── ast/         # Abstract Syntax Tree nodes implementing position-aware interface
├── parser/      # Recursive descent Pratt parser with optional semicolons/types
├── object/      # Runtime value representation and Environment scope tracking
├── evaluator/   # AST Tree-walker evaluator and built-in stdlib registration
├── repl/        # Interactive REPL with multiline block detection
├── samples/     # Example RustX source code files (.rx)
├── main.go      # CLI entry point (run, build, repl, version)
└── README.md
```

---

## 🚀 Getting Started & Building

### Prerequisites
- Go 1.20+

### Building the CLI Binary
To compile the `rustx` executable in the root directory:

```bash
go build -o rustx main.go
```

---

## 💻 CLI Commands

### 1. Execute a Script
Run a RustX script directly:
```bash
./rustx run samples/hello.rx
```
*Shorthand execution is also supported:*
```bash
./rustx samples/hello.rx
```

### 2. Interactive REPL
Start the interactive REPL:
```bash
./rustx repl
# or simply:
./rustx
```
Example session:
```text
RustX Starter v0.1.0 Interactive REPL
Type statements or expressions to evaluate. Type 'exit' to quit.

>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
>>> exit
```

### 3. Syntax Validation / Build Check
Check a `.rx` file for syntax errors without executing it:
```bash
./rustx build samples/demo.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## 📝 Example Source Code

### Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello RustX Starter")
}
```

### Recursion & Math (`samples/factorial.rx`)
```rustx
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

### Control Flow & Standard Library (`samples/demo.rx`)
```rustx
func main() {
    var sum = 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    println("Sum 1..10:", sum)

    msg := "apple,banana,cherry"
    parts := split(msg, ",")
    println("Upper joined:", to_upper(join(parts, " - ")))
}
```

---

## 🧪 Testing

Run all unit tests across all packages (`lexer`, `parser`, `evaluator`):

```bash
go test ./...
```
