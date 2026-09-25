# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by **Rust** and **Go**. Designed for simplicity, speed, readability, and seamless future extensibility, RustX Starter blends Rust's clean syntax with Go's pragmatic usability.

---

## Key Features

- **Blended Syntax:** Inspired by Go and Rust (`func`, `var`, `:=`, optional semicolons, type annotations).
- **Core Data Types:** Integers (`int`), Floats (`float`), Booleans (`bool`), Strings (`string`), and Arrays (`[elem1, elem2]`).
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), Relational (`==`, `!=`, `<`, `>`, `<=`, `>=`), Postfix (`++`, `--`), and Logical (`&&`, `||`).
- **Control Flow:** `if / else`, `while`, and flexible Go-style `for` loops (`for init; cond; post { ... }` or `for cond { ... }`).
- **Functions:** Parameter type annotations, return types, recursive call support, and top-level entry point `func main()`.
- **Built-in Standard Library:**
  - **I/O:** `print()`, `println()`
  - **Strings & Collections:** `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math:** `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time:** `now()`, `sleep()`
  - **Filesystem:** `read_file()`, `write_file()`
- **Comments:** Single-line (`//`) and multi-line (`/* ... */`).
- **Interactive REPL:** Multiline input support with persistent evaluation environment.
- **Human-Readable Error Handling:** Position-aware diagnostics (`SyntaxError`, `RuntimeError`).

---

## Architecture

RustX Starter is designed with a clean, modular architecture located in separate packages:

```
├── token/        # Token definitions, keywords, operator maps, position info
├── errors/       # Standardized human-readable SyntaxError & RuntimeError formatting
├── lexer/        # Hand-written Lexer/Tokenizer tracking line/col & handling escapes/comments
├── ast/          # Extensible AST nodes (Statements, Expressions, Function Literals, etc.)
├── parser/       # Recursive descent + Pratt Parser with configurable precedences
├── object/       # Runtime Value Objects & Enclosed Environment scoping system
├── evaluator/    # Evaluator engine & Standard Library built-in functions
├── repl/         # Interactive REPL prompt with multiline input tracking
├── samples/      # Sample RustX scripts (.rx)
└── main.go       # CLI main entry point
```

### Future Extensibility
The architecture intentionally defers borrow checking, lifetimes, and full ownership models. However, all AST nodes contain token location data (`GetToken()`) and Environment scopes are enclosed, allowing future phases to add:
- Ownership & Borrow Checker static pass
- Structs, Enums, Traits/Interfaces
- Concurrency & Async/Await
- Modules & Imports
- Bytecode VM / Native compilation backend

---

## Installation & Build Instructions

### Prerequisites
- Go 1.20+ installed.

### Building RustX Binary
From the root directory:
```bash
go build -o rustx main.go
```

### Running Tests
To run unit tests across all packages:
```bash
go test ./...
```

---

## CLI Usage

### 1. Run Script Directly
```bash
./rustx run samples/hello.rx
```
Or shorthand execution:
```bash
./rustx samples/factorial.rx
```

### 2. Interactive REPL
```bash
./rustx repl
```
Or simply:
```bash
./rustx
```

### 3. Check Version
```bash
./rustx version
```

### 4. Build / Syntax Validation
Validates syntax without executing the script:
```bash
./rustx build samples/features.rx
```

---

## Code Examples

### Example 1: Hello World
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Example 2: Recursive Factorial
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### Example 3: Control Flow & Standard Library
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

for i := 0; i < 5; i++ {
    println("i =", i)
}

words := split("rust,go,rustx", ",")
println("Upper joined:", to_upper(join(words, " - ")))
```

---

## License

MIT License.
