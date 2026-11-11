# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="220" />
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**, focused on simplicity, speed, readability, and future extensibility. It bridges the ergonomics of Go with the expressiveness and foundational architecture of Rust.

> **Note on Architecture:**
> RustX Starter is built as the solid architectural foundation for a future systems language. It intentionally omits borrow checker, lifetimes, and complex ownership rules at this phase, behaving as a lightweight scripting + systems hybrid while structuring its AST and parser for zero-rewrite future enhancements.

---

## Features

- **Blended Go & Rust Syntax**:
  - `func` for function declarations
  - `var` for mutable variable declarations
  - `:=` shorthand declaration and assignment
  - Optional semicolons and optional condition parentheses
- **Data Types**:
  - `int` (64-bit integer)
  - `float` (64-bit float)
  - `bool` (`true`, `false`)
  - `string` (with escape sequences like `\n`, `\t`, `\"`, `\\`)
  - `array` (lists and indexing `arr[i]`, with index assignment `arr[i] = val`)
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!` (with short-circuit evaluation)
  - Postfix: `++`, `--`
- **Control Flow**:
  - `if`, `else if`, `else`
  - `while` loops
  - Go-style 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`)
- **Functions**:
  - First-class and named functions
  - Type annotations (e.g. `func add(a int, b int) int`)
  - Full recursion support
  - Automatic `main()` entrypoint invocation when running scripts
- **Comments**:
  - `//` single-line comments
  - `/* ... */` multi-line comments
- **Standard Library**:
  - **I/O & Printing**: `print()`, `println()`
  - **Strings**: `len()`, `to_upper()`, `to_lower()`, `contains()`, `split()`, `trim()`
  - **Math**: `abs()`, `sqrt()`, `pow()`, `min()`, `max()`
  - **Time**: `time_now()`, `sleep_ms()`
  - **Filesystem**: `read_file()`, `write_file()`
  - **Collections**: `push()`
- **Developer Experience**:
  - Human-readable `SyntaxError` and `RuntimeError` messages with line and column numbers.
  - Interactive multiline REPL with brace and parenthesis tracking.
  - CLI supporting `run`, `build` (syntax checking), `version`, and `repl`.

---

## Project Structure

```text
rustx/
├── ast/            # Abstract Syntax Tree definitions & extensible node models
├── errors/         # Human-readable syntax & runtime error diagnostics
├── evaluator/      # Runtime evaluation engine & standard library builtins
├── lexer/          # Lexical scanner & tokenizer
├── parser/         # Pratt & recursive-descent parser
├── object/         # Object system, environment, & lexical scoping
├── repl/           # Multiline interactive REPL prompt
├── samples/        # Example RustX programs (.rx)
├── token/          # Token constants & keyword definitions
├── main.go         # CLI entry point
└── README.md       # Project documentation
```

---

## Installation & Building

RustX is written in Go and requires Go 1.20+ (tested on Go 1.24).

### Build CLI Binary
```bash
go build -o rustx main.go
```

### Run Tests
```bash
go test ./...
```

---

## Usage & CLI

### 1. Interactive REPL
Launch the REPL by running `rustx` or `rustx repl`:
```bash
$ ./rustx
RustX Starter v0.1
Type expressions or statements. Press Ctrl+D (or type 'exit') to quit.
>>> x := 5
>>> y := 10
>>> x + y
15
>>> func greet(name) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### 2. Run a Script
Run `.rx` files using `rustx run <file>` or shorthand `rustx <file>`:
```bash
./rustx run samples/01_hello.rx
# or
./rustx samples/01_hello.rx
```

### 3. Build / Validate Syntax
Check for syntax errors without executing:
```bash
./rustx build samples/02_factorial.rx
```

### 4. Check Version
```bash
./rustx version
```

---

## Example Code

### Hello World & Main Entrypoint
```rustx
func main() {
    println("RustX Starter")
}
```

### Recursive Factorial
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### Variables & Conditionals
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
} else {
    println("Unknown language")
}
```

### Loops
```rustx
// Go-style for loop
for i := 0; i < 5; i++ {
    println("i =", i)
}

// While loop
count := 3
while count > 0 {
    println("count =", count)
    count--
}
```

### Standard Library Basics
```rustx
// String manipulation
s := "  rustx starter  "
println(to_upper(trim(s)))

// Math
println("Square root of 144:", sqrt(144))
println("2^10:", pow(2, 10))

// File I/O
write_file("greeting.txt", "Hello World from RustX!")
println(read_file("greeting.txt"))
```

---

## Future Roadmap

RustX Starter is designed so that future systems capabilities can be layered on without redesigning the core AST or lexer:
- **Phase 2**: Structs, enums, pattern matching, methods.
- **Phase 3**: Ownership semantics, move vs. copy, reference types.
- **Phase 4**: Borrow checker and non-lexical lifetimes.
- **Phase 5**: Modules, package manager, and bytecode VM / native compiler (LLVM or Cranelift).
