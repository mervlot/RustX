# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a minimal, clean, modular programming language that blends syntax inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and future extensibility, RustX Starter serves as a lightweight scripting + systems hybrid language foundation.

> **Note on Architecture:** RustX Starter currently operates as a lightweight interpreted scripting language without an active borrow checker or lifetime annotations. The AST, evaluator, environment, and parser are structured modularly to allow adding ownership, borrow checking, and a VM or native compilation backend in future releases without rewriting the core framework.

---

## Language Features

- **Blended Go & Rust Syntax:**
  - `func` for function declarations
  - `var` for mutable variable definitions
  - `:=` shorthand assignment operator
  - Optional static type annotations (e.g., `func add(a int, b int) int`)
  - Curly brace blocks with optional semicolons
- **Data Types:**
  - `int` (64-bit integer)
  - `float` (64-bit float)
  - `bool` (`true`, `false`)
  - `string` (with full backslash escape sequence decoding)
  - `array` (`[1, 2, 3]`)
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||` (with short-circuit evaluation)
  - Postfix Increment / Decrement: `++`, `--`
- **Control Flow:**
  - `if`, `else if`, `else`
  - `while` loops
  - `for` loops (supports both 3-part `for i := 0; i < 10; i++` and condition-only `for i < 10`)
  - Optional parentheses around conditions and loop headers
- **Comments:**
  - Single-line: `// ...`
  - Multi-line: `/* ... */`
- **Standard Library Built-ins:**
  - Printing: `print`, `println`
  - String operations: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math utilities: `abs`, `pow`, `sin`, `cos`, `tan`
  - System / Time: `sleep`, `now`
  - Filesystem basics: `read_file`, `write_file`

---

## Project Architecture

```text
rustx/
├── token/      # Token types, definitions, and keyword map
├── errors/     # Structured SyntaxError and RuntimeError diagnostics
├── lexer/      # Tokenizer/lexer supporting comments, escapes, & numbers
├── ast/        # Position-aware AST nodes for statements & expressions
├── parser/     # Pratt / Recursive Descent parser
├── object/     # Runtime object system and lexical Environment scopes
├── evaluator/  # AST interpreter and built-in function registry
├── repl/       # Interactive prompt with multiline block tracking
├── samples/    # Sample .rx programs demonstrating language features
├── main.go     # CLI entry point (run, repl, version, build, direct exec)
└── README.md
```

---

## Installation & Usage

### Prerequisites
- [Go 1.20+](https://golang.org/dl/)

### Building from Source

```bash
go build -o rustx main.go
```

### CLI Commands

#### 1. Interactive REPL Prompt
```bash
./rustx repl
# or simply:
./rustx
```
*Sample REPL output:*
```text
RustX Starter v0.1
Type 'exit' to quit.
>>> x := 5
>>> println(x * 2)
10
```

#### 2. Running a RustX Script
```bash
./rustx run samples/hello.rx
# or direct file shorthand:
./rustx samples/hello.rx
```

#### 3. Syntax Verification / Build Check
Validate syntax without executing runtime side effects:
```bash
./rustx build samples/factorial.rx
```

#### 4. Version Check
```bash
./rustx version
```

---

## Code Examples

### 1. Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Recursive Factorial (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5:")
    println(factorial(5))
}
```

### 3. Control Flow & Standard Library (`samples/features.rx`)
```rustx
func main() {
    println("--- Control Flow & Loops ---")
    var sum = 0
    for i := 1; i <= 5; i++ {
        sum = sum + i
        println(i)
    }
    println("Sum of 1 to 5:")
    println(sum)

    println("--- Standard Library ---")
    msg := "RustX Starter"
    println(to_upper(msg))
    println(contains(msg, "Rust"))
    println(pow(2, 3))
}
```

---

## Running Tests

Run all unit tests across the repository:

```bash
go test ./...
```

---

## Future Roadmap

- [ ] Structs and Enums
- [ ] Traits / Interfaces
- [ ] Concurrency (`go` / `async` routines)
- [ ] Ownership & Borrow Checker Static Analysis Phase
- [ ] Bytecode VM compiler and Native Backend
