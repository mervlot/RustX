# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by **Rust** and **Go**. It blends Rust's clean modern syntax with Go's simplicity and execution speed.

> **Note on Architecture & Future Readiness:**
> The current version serves as the foundation phase. It intentionally omits ownership, borrow checking, and lifetimes while maintaining a clean, decoupled modular architecture (`lexer`, `parser`, `ast`, `object`, `evaluator`, `errors`, `repl`) designed so advanced memory safety models, bytecode VMs, and native compilers can be attached without rewriting the frontend.

---

## Features Overview

- **Syntactic Blend:** `func` for functions, `var` for declarations, `:=` for shorthand, strong dynamic types, curly braces, and optional semicolons.
- **Data Types:** `int`, `float`, `bool`, `string`, `Array`.
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
  - Postfix: `++`, `--`
- **Control Flow:** `if` / `else if` / `else`, `while` loops, and 3-part or condition-only `for` loops.
- **Standard Library Built-ins:**
  - I/O: `print`, `println`
  - Strings & Collections: `len`, `split`, `join`, `to_upper`, `to_lower`, `contains`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time & System: `sleep`, `now`, `read_file`, `write_file`
- **CLI & Interactive REPL:** Single-file script runner, multiline interactive REPL, syntax verification build mode, and shorthand execution.

---

## Directory Structure

```text
.
├── ast/            # AST node definitions with position metadata
├── errors/         # Human-readable SyntaxError and RuntimeError formatters
├── evaluator/      # AST Evaluator and Standard Library built-in functions
├── lexer/          # Tokenizer with line & column tracking and comment support
├── object/         # Runtime representation and lexical environment scopes
├── parser/         # Pratt Recursive Descent Parser
├── repl/           # Multiline REPL interface
├── samples/        # Demonstration scripts (.rx)
├── token/          # Token definitions and keyword lookup
├── main.go         # CLI entrypoint
└── README.md
```

---

## Quick Start & Installation

### Requirements
- [Go](https://go.dev/) 1.20 or newer.

### Building RustX CLI
Build the CLI binary locally:

```bash
go build -o rustx main.go
```

Or run directly using Go:

```bash
go run main.go [command]
```

---

## CLI Usage

### Running Scripts
Execute any `.rx` file:

```bash
./rustx run samples/hello.rx
```

Or use shorthand execution:

```bash
./rustx samples/factorial.rx
```

### Interactive REPL
Start the interactive prompt:

```bash
./rustx repl
```

Example session:
```text
RustX Starter v0.1
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func greet(name) {
...     return "Hello " + name
... }
>>> greet("RustX")
"Hello RustX"
```

### Syntax Verification (Build Check)
Check syntax without running the script:

```bash
./rustx build samples/features.rx
```

### Version
Check version:

```bash
./rustx version
```

---

## Syntax Examples

### 1. Variables & Functions

```rustx
var age int = 25;
name := "Muhammad";

func add(a int, b int) int {
    return a + b;
}

println("Name:", name, "Age:", age)
println("5 + 10 =", add(5, 10))
```

### 2. Control Flow & Recursion

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

### 3. Loops

```rustx
func main() {
    // 3-part for loop
    sum := 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    println("Sum of 1..10:", sum)

    // While loop
    count := 3
    while count > 0 {
        println("Countdown:", count)
        count--
    }
}
```

---

## Standard Library Quick Reference

| Category | Function | Example |
| :--- | :--- | :--- |
| **Console** | `print(...)`, `println(...)` | `println("Hello World")` |
| **Strings** | `len(str)`, `split(str, sep)`, `join(arr, sep)` | `split("a,b", ",")` |
| **Strings** | `to_upper(str)`, `to_lower(str)`, `contains(str, sub)` | `to_upper("rustx")` |
| **Math** | `abs(x)`, `pow(base, exp)` | `pow(2, 8)` |
| **Math** | `sin(x)`, `cos(x)`, `tan(x)` | `sin(3.14159 / 2)` |
| **Time** | `sleep(ms)`, `now()` | `sleep(100)` |
| **Filesystem** | `read_file(path)`, `write_file(path, data)` | `write_file("data.txt", "content")` |

---

## Error Diagnostics

Errors provide line and column information:

```text
SyntaxError: expected next token to be ), got IDENT instead at line 5 column 12
RuntimeError: division by zero at line 10 column 4
```

---

## Running Unit Tests

Run all unit tests across the codebase:

```bash
go test ./...
```
