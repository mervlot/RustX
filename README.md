# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. Designed with simplicity, speed, readability, and future extensibility in mind, RustX Starter blends the best ergonomics of Rust and Go into a lightweight scripting and systems hybrid language.

> **Note on Architecture & Future Compatibility:**
> The current version (v0.1) serves as the foundational phase. Advanced memory safety systems (borrow checking, ownership models, lifetimes) are deliberately omitted for now, but the underlying AST and compiler/interpreter architecture are structured to support them, alongside structs, enums, traits, async/await, and native compilation in future iterations without architectural rewrites.

---

## Language Features

- **Blended Rust & Go Syntax**: `func` for function declarations, `var` and `:=` shorthand for variable bindings, curly brace blocks, optional semicolons, and unparenthesized control flows.
- **Dynamic Primitive Types**: `int`, `float`, `bool`, and `string`.
- **Control Flow**: `if` / `else if` / `else`, `while` loops, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Relational (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Prefix (`!`, `-`), Postfix (`++`, `--`).
- **Functions**: Named and anonymous functions, type annotations, return values, recursive calls.
- **Built-in Standard Library**:
  - Printing: `print()`, `println()`
  - String Operations: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Math Utilities: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - Time & Filesystem: `now()`, `sleep()`, `read_file()`, `write_file()`
- **Human-Readable Errors**: Precise `SyntaxError` and `RuntimeError` reporting with line and column numbers.

---

## Directory Structure

```text
rustx/
├── ast/         # AST Node definitions (implementing Position-aware interfaces)
├── errors/      # Custom SyntaxError and RuntimeError diagnostics
├── evaluator/   # Interpreter engine, scope management, and built-ins
├── lexer/       # Tokenizer with line/column tracking & comment support
├── object/      # Runtime object system and Environment scopes
├── parser/      # Pratt recursive descent parser
├── repl/        # Interactive stateful REPL with multiline support
├── samples/     # Sample RustX Starter scripts (.rx)
├── token/       # Token types and keyword mappings
├── main.go      # CLI entry point
└── go.mod       # Go module declaration
```

---

## Getting Started

### Prerequisites
- [Go 1.20+](https://golang.org/dl/)

### Installation & Building

Clone the repository and build the binary:

```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

RustX Starter supports multiple execution modes:

### 1. Run a Script
Run `.rx` files directly using the `run` command or shorthand execution:

```bash
./rustx run samples/hello.rx
# or
./rustx samples/hello.rx
```

### 2. Interactive REPL
Launch the stateful multiline REPL prompt:

```bash
./rustx repl
# or simply
./rustx
```

REPL Example:
```text
RustX Starter v0.1 Prompt
Type 'exit' or Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(x, y))
30
```

### 3. Check Version
```bash
./rustx version
```

### 4. Build / Validate Code
Check syntax and lexical correctness of a file without executing it:

```bash
./rustx build samples/factorial.rx
```

---

## Code Examples

### Example 1: Hello World
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### Example 2: Recursion (Factorial)
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

### Example 3: Control Flow & Loops
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

var sum = 0
for i := 0; i < 5; i++ {
    sum = sum + i
}
println("Sum from 0 to 4 is:", sum)
```

---

## License

MIT
