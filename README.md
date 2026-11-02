<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="240" />
</p>

# RustX Starter

**RustX Starter** is a minimal, clean, and extensible programming language inspired by **Rust** and **Go**. It balances modern syntactic clarity, rapid execution, readability, and a decoupled architecture prepared for future systems capabilities.

> **Note on Architecture:** Advanced memory safety systems (such as borrow checking, explicit lifetimes, and ownership tracking) are intentionally designed as future phases. The foundation uses a modular AST and interpreter architecture that allows these systems and bytecode/native compilation passes to be introduced incrementally without architectural rewrites.

---

## Features

- **Blended Rust + Go Syntax**:
  - `func` keyword for functions with optional parameter and return types
  - `var` declarations and `:=` shorthand definitions
  - Block scoping with `{ ... }`
  - Optional semicolons
  - Post-increment and post-decrement (`++`, `--`)
  - Arithmetic (`+`, `-`, `*`, `/`, `%`) and logical short-circuiting operators (`&&`, `||`, `!`)
  - Comparison operators (`==`, `!=`, `<`, `>`, `<=`, `>=`)
- **Data Types**:
  - Integers (`int`)
  - Floating-point numbers (`float`)
  - Booleans (`bool`: `true`, `false`)
  - Strings (`string` with escape sequences `\n`, `\t`, `\"`, `\\`)
  - Arrays (`[1, 2, 3]`)
- **Control Flow**:
  - `if` / `else if` / `else`
  - `while` loops
  - Go-style 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`) and condition-only `for` loops
  - Parentheses optional around conditions and loop headers
- **Functions**:
  - First-class and named functions
  - Closures and lexical scoping
  - Full recursion support
  - Optional entry point: `func main()` automatically invoked when executing scripts
- **Interactive REPL**:
  - Multi-line block detection
  - Persistent lexical environment
- **Standard Library Built-ins**:
  - **Console**: `print(...)`, `println(...)`
  - **Collections & Strings**: `len(x)`, `str_len(s)`, `str_split(s, sep)`, `str_contains(s, sub)`, `str_upper(s)`, `str_lower(s)`, `str_trim(s)`
  - **Math**: `math_sqrt(x)`, `math_abs(x)`, `math_pow(base, exp)`, `math_min(a, b)`, `math_max(a, b)`
  - **Time**: `time_now()`, `time_sleep(ms)`
  - **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Human-Readable Diagnostics**:
  - Syntax error reporting with line and column positions
  - Runtime errors that gracefully report issues without runtime panics

---

## Architecture

The project is structured into modular decoupled packages:

```
├── ast/          # Abstract Syntax Tree node definitions
├── errors/       # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/    # AST interpreter, scoping, and built-in standard library
├── lexer/        # Tokenizer with position tracking (line & column)
├── object/       # Runtime object representation and Environment scope chain
├── parser/       # Pratt / recursive descent parser with precedence engine
├── repl/         # Interactive REPL supporting multiline expressions
├── samples/      # Example .rx programs
├── token/        # Token definitions and keyword lookup
└── main.go       # CLI tool and entry point
```

---

## Installation & Building

RustX Starter requires **Go 1.20+**.

### Build Binary
```bash
go build -o rustx main.go
```

### Run Tests
```bash
go test -v ./...
```

---

## CLI Usage

### 1. Interactive REPL
```bash
./rustx repl
# or simply:
./rustx
```

Example session:
```text
RustX Starter v0.1
Type 'exit' or Ctrl+D to quit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> exit
```

### 2. Execute a Script
```bash
./rustx run samples/01_hello.rx
# or using shorthand:
./rustx samples/01_hello.rx
```

### 3. Syntax Verification / Build
```bash
./rustx build samples/02_factorial.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## Language Examples

### Example 1: Hello World
```rust
func main() {
    println("Hello World from RustX!")
}
```

### Example 2: Recursive Factorial
```rust
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

### Example 3: Loops & Conditionals
```rust
func main() {
    for i := 0; i < 5; i++ {
        println("Iteration:", i)
    }

    count := 3
    while count > 0 {
        println("Countdown:", count)
        count--
    }
}
```

### Example 4: Standard Library
```rust
func main() {
    greeting := "  hello rustx  "
    cleaned := str_upper(str_trim(greeting))
    println(cleaned) // HELLO RUSTX

    println("Sqrt of 64:", math_sqrt(64))
    println("Current time:", time_now())

    write_file("output.txt", "RustX works!")
    println("Read back:", read_file("output.txt"))
}
```
