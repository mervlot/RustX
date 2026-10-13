# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a lightweight, compiled/interpreted programming language inspired by **Rust** and **Go**. It blends the simplicity, clean syntax, and fast ergonomics of Go with the power and future extensibility of Rust.

The current foundation phase operates as a hybrid scripting/systems language, designed cleanly so future versions can integrate ownership models, borrow checkers, lifetimes, structs, enums, modules, and native/bytecode compilation without rewriting the core architecture.

---

## Language Features

- **Blended Go/Rust Syntax**:
  - `func` for functions
  - `var` for mutable variable declarations
  - `:=` shorthand assignment operator
  - Optional type annotations in function signatures (e.g. `func add(a int, b int) int`)
  - Optional semicolons and unparenthesized control flow conditions
- **Data Types**: `int`, `float`, `bool`, `string`, `array`, `null`
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Post-increment/decrement: `++`, `--`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
- **Control Flow**:
  - `if`, `else if`, `else`
  - `while` loops
  - 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`)
- **Functions & Recursion**: First-class functions and recursive calls
- **Standard Library Built-ins**:
  - Printing: `println()`, `print()`
  - Collections & Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - System & Filesystem: `now()`, `sleep()`, `read_file()`, `write_file()`
- **Comments**:
  - `//` single-line comments
  - `/* ... */` multi-line comments

---

## Directory & Architecture Overview

The codebase is organized into clean, modular Go packages:

```text
rustx/
├── token/      # Token types, positions, and keyword definitions
├── errors/     # Structured SyntaxError and RuntimeError diagnostics with line/col info
├── lexer/      # Lexical analyzer supporting comments, escape sequences, and position tracking
├── ast/        # Abstract Syntax Tree nodes for statements, expressions, and literals
├── parser/     # Pratt recursive descent parser
├── object/     # Runtime values, object representations, and Environment scoping
├── evaluator/  # Core tree-walk evaluator and standard library built-ins
├── repl/       # Interactive REPL prompt with multi-line brace-balancing support
├── samples/    # Sample programs demonstrating RustX language features
├── main.go     # CLI entry point
└── README.md   # Documentation
```

---

## Quick Start & Building

### Prerequisites
- **Go 1.20+** installed

### Build the Executable

```bash
go build -o rustx main.go
```

### Running Tests

```bash
go test ./...
```

---

## CLI Usage

### Direct File Execution (Shorthand)
```bash
./rustx samples/01_hello.rx
```

### Run Command
```bash
./rustx run samples/02_factorial.rx
```

### Check Syntax / Validate Script (Build)
```bash
./rustx build samples/03_loops_and_arrays.rx
```

### Interactive REPL
```bash
./rustx repl
# or simply run without arguments:
./rustx
```

Example REPL session:
```text
RustX Starter v0.1
Type "exit" or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
```

### Version
```bash
./rustx version
```

---

## Sample Code Examples

### 1. Hello World
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Recursive Factorial
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

### 3. Loops and Arrays
```rustx
func main() {
    numbers := [10, 20, 30, 40, 50]
    sum := 0

    for i := 0; i < len(numbers); i++ {
        println("Index", i, "value:", numbers[i])
        sum = sum + numbers[i]
    }

    println("Total sum:", sum)
}
```

---

## Future Architecture & Scalability

RustX Starter is explicitly architected to allow future language evolutions:
1. **Ownership & Borrow Checking**: AST nodes and environment scopes store source tracking metadata so ownership checkers can pass over the AST prior to runtime evaluation.
2. **Bytecode VM or Native Compilation**: The clean separation between `parser` -> `ast` -> `evaluator` makes replacing the evaluator with an IR/bytecode compiler straightforward.
3. **Structs, Enums, and Interfaces**: AST nodes extend `ast.Statement` and `ast.Expression` interfaces without breaking existing grammar rules.
