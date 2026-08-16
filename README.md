# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a lightweight, compiled/interpreted scripting + systems hybrid language inspired by **Rust** and **Go**. It combines Rust's modern syntax and safety-oriented design goals with Go's simplicity, speed, and clean readability.

---

## 🌟 Key Features

- **Blended Syntax**: Utilizes `func` for functions, `var` for mutable variable declarations, and `:=` for shorthand declaration.
- **Data Types**: Support for `int`, `float`, `bool`, `string`, and `array`.
- **Control Flow**:
  - `if` / `else if` / `else` clauses
  - `while` loops
  - Go-style 3-part `for` loops (e.g. `for i := 0; i < 10; i++`)
- **Functions & Recursion**: Named and anonymous function literals with optional type annotations and recursive invocation.
- **Rich Operator Support**: Arithmetic (`+`, `-`, `*`, `/`), logical (`&&`, `||` with short-circuiting), postfix (`++`, `--`), and indexing (`[index]`).
- **Standard Library**:
  - **Printing**: `print()`, `println()`
  - **Strings & Collections**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time & IO**: `sleep()`, `now()`, `read_file()`, `write_file()`
- **Precise Error Diagnostics**: Human-readable syntax errors and runtime errors with exact line and column numbers.
- **Interactive REPL**: Interactive prompt supporting multiline input tracking.

---

## 📁 Architecture Overview

RustX Starter is designed with a modular, scalable architecture to enable future expansion (e.g. adding borrow checkers, structs, generics, or native compilation):

```
.
├── ast/        # Abstract Syntax Tree nodes
├── errors/     # Formatted diagnostic errors (SyntaxError, RuntimeError)
├── evaluator/  # Tree-walking interpreter & standard library built-ins
├── lexer/      # Lexical analyzer / Tokenizer
├── object/     # Runtime object system and Environment (Lexical Scoping)
├── parser/     # Recursive Descent + Pratt Expression Parser
├── repl/       # Interactive REPL console
├── samples/    # Example .rx programs
├── token/      # Token type definitions and keyword maps
└── main.go     # CLI entry point
```

---

## 🚀 Quick Start & CLI Usage

### Build CLI Binary
```bash
go build -o rustx main.go
```

### Run a Script
```bash
./rustx run samples/01_hello.rx
# OR shorthand execution:
./rustx samples/01_hello.rx
```

### Validate Script Syntax
```bash
./rustx build samples/02_factorial.rx
```

### Interactive REPL
```bash
./rustx repl
# OR run without arguments:
./rustx
```

In the REPL:
```rust
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(15, 25)
40
```

---

## 📝 Code Examples

### Hello World
```rust
func main() {
    println("Hello World from RustX!")
}
```

### Recursive Factorial
```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### Loops & Conditionals
```rust
var sum = 0
for i := 1; i <= 5; i++ {
    sum = sum + i
}

if sum > 10 {
    println("Sum is greater than 10:", sum)
}
```

---

## 🧪 Testing

Run all unit tests across the codebase:

```bash
go test ./...
```
