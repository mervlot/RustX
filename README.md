# RustX Starter

![RustX Starter Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by Rust and Go, designed for simplicity, speed, readability, and future extensibility.

---

## 🌟 Key Features

- **Blended Syntax:** Inspired by Go and Rust (`func`, `var`, `:=`, optional semicolons, curly brace blocks).
- **Core Types:** Integers (`int`), Floats (`float`), Booleans (`bool`), Strings (`string`), and Arrays (`[...]`).
- **Control Flow:** `if` / `else`, `while` loops, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Functions:** Named functions, anonymous functions, parameter typing, return type annotations, and recursion.
- **Built-in Standard Library:** Built-in utilities for I/O (`println`, `print`), strings (`len`, `str_upper`, `str_lower`, `str_contains`), math (`abs`, `pow`, `sqrt`), time (`time_now`, `sleep`), and filesystem operations (`read_file`, `write_file`, `file_exists`).
- **Interactive REPL:** Interactive REPL environment with multiline input tracking and block balance awareness.
- **Human-Readable Diagnostics:** Structured syntax and runtime diagnostics reporting exact line and column locations without interpreter panics.

---

## 🏗 System Architecture & Directory Structure

```
.
├── ast/         # Abstract Syntax Tree (AST) node definitions
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walking interpreter and standard library built-ins
├── lexer/       # Tokenizer scanning source characters into tokens
├── object/      # Runtime object types and Environment lexical scoping
├── parser/      # Pratt recursive descent parser
├── repl/        # Interactive REPL session handler
├── samples/     # Demonstrative .rx example scripts
├── token/       # Token definitions and keyword mapping
├── main.go      # CLI entry point
└── README.md    # Language documentation
```

### Future Architecture Readiness
The AST nodes and environment hierarchy are designed to easily extend toward memory management models, ownership checking, modules/imports, async/await, structs, enums, and native byte-code generation.

---

## 🚀 Quick Start & CLI Usage

### Build CLI Binary
```bash
go build -o rustx main.go
```

### Run a Script
```bash
./rustx run samples/hello.rx
# or directly:
./rustx samples/hello.rx
```

### Interactive REPL Prompt
```bash
./rustx repl
```
Example:
```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> func double(n int) int { return n * 2 }
>>> println(double(x))
20
```

### Syntax Validation / Build Check
```bash
./rustx build samples/factorial.rx
```

### Display Version
```bash
./rustx version
```

---

## 💻 Syntax & Code Examples

### 1. Hello World
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Variables & Operators
```rustx
var age int = 16
name := "Muhammad"
is_ready := true

count := 0
count++
println(name, "age:", age, "count:", count)
```

### 3. Recursive Functions
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### 4. Loops & Control Flow
```rustx
// For Loop
var sum = 0
for i := 1; i <= 5; i++ {
    sum = sum + i
}
println("Sum:", sum)

// While Loop
n := 3
while n > 0 {
    println("Countdown:", n)
    n--
}
```

---

## 🧪 Testing

Run all unit tests across the repository:
```bash
go test ./...
```
