# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a minimal, clean, and scalable scripting + systems hybrid programming language designed as a blend of **Rust** and **Go**. It balances modern, clean syntax, simplicity, high performance execution, and human-readable diagnostics while laying an architectural foundation for future language features such as static types, borrow checking, ownership models, modules, and native compilation.

---

## Features

- **Blended Rust/Go Syntax**: `func` for functions, `var` for declarations, `:=` for shorthand definitions, curly brace blocks, optional semicolons.
- **Dynamic & Strong Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Flexible Control Flow**: `if` / `else if` / `else`, `while` loops, single-condition and 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), Postfix (`++`, `--`).
- **Rich Built-in Standard Library**:
  - **I/O & Console**: `print(...)`, `println(...)`
  - **Math**: `abs(x)`, `pow(base, exp)`, `sin(x)`, `cos(x)`, `tan(x)`
  - **Strings & Collections**: `len(val)`, `contains(s, sub)`, `split(s, sep)`, `join(arr, sep)`, `to_upper(s)`, `to_lower(s)`
  - **Time**: `now()`, `sleep(ms)`
  - **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Interactive Multiline REPL**: Preserves state across evaluations and automatically waits for balanced block braces `{}` and parentheses `()`.
- **Human-Readable Diagnostics**: Precise error messages containing line and column position.

---

## Architecture Overview

RustX Starter follows a modular, decoupled compiler/interpreter pipeline:

```
[ Source Code (.rx) ]
          │
          ▼
    ┌───────────┐
    │   Lexer   │ Token Stream with Line & Column Locations
    └─────┬─────┘
          │
          ▼
    ┌───────────┐
    │  Parser   │ Pratt Expression Parser + Recursive Descent
    └─────┬─────┘
          │
          ▼
    ┌───────────┐
    │    AST    │ Abstract Syntax Tree Nodes
    └─────┬─────┘
          │
          ▼
    ┌───────────┐
    │ Evaluator │ Environment Scoping, Builtins & Execution Engine
    └───────────┘
```

### Directory Structure

```
.
├── ast/         # AST interface and concrete expression/statement definitions
├── errors/      # Standardized human-readable error formatting
├── evaluator/   # Tree-walk evaluator & Standard library builtins
├── lexer/       # Tokenizer with line & column tracking and comment stripping
├── object/      # Runtime objects & Environment scope hierarchy
├── parser/      # Pratt expression and recursive descent statement parser
├── repl/        # Multiline interactive REPL prompt engine
├── samples/     # Example .rx scripts
├── token/       # Token types and keywords table
├── main.go      # CLI entry point
└── go.mod       # Go module file
```

---

## Installation & Build Instructions

### Requirements
- [Go](https://golang.org/) v1.20 or newer

### Building from Source

```bash
# Clone repository
git clone https://github.com/rustx/rustx.git
cd rustx

# Build the rustx binary
go build -o rustx main.go
```

---

## Usage

### 1. Execute a Script
```bash
./rustx run samples/01_hello.rx
# OR run directly:
./rustx samples/01_hello.rx
```

### 2. Interactive REPL
```bash
./rustx repl
```
Example session:
```
RustX Starter REPL v0.1.0
Type 'exit' or press Ctrl+D to quit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func greet(name) {
...     return "Hello " + name
... }
>>> println(greet("World"))
Hello World
```

### 3. Check Syntax (Build Command)
```bash
./rustx build samples/02_factorial.rx
```

### 4. Version
```bash
./rustx version
```

---

## Language Syntax Examples

### Functions & Recursion
```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    result := factorial(5)
    println("Factorial of 5 is:", result)
}
```

### Control Flow & Arrays
```rust
func main() {
    var sum = 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    println("Sum =", sum)

    arr := [10, 20, 30]
    for idx := 0; idx < len(arr); idx++ {
        println("arr[", idx, "] =", arr[idx])
    }
}
```

---

## Running Tests

Run the full unit test suite across all packages:

```bash
go test ./...
```
