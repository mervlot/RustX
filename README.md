# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a minimal, clean, scalable programming language inspired by **Rust** and **Go**. It balances simplicity, performance, readability, and future extensibility.

> **Architecture Note:** RustX Starter is designed as a foundational language engine. Ownership models, borrow checkers, lifetimes, and advanced concurrency primitives are deliberately excluded from this foundation phase, but the modular AST and environment architecture allow these systems to be added in future iterations without rewriting the compiler core.

---

## Language Overview

- **File Extension:** `.rx`
- **Keywords:** `func`, `var`, `if`, `else`, `while`, `for`, `return`, `true`, `false`
- **Shorthand Variable Assignment:** `:=`
- **Types:** `int`, `float`, `bool`, `string`, `array` (type annotations in function signatures and variable declarations are supported)
- **Control Flow:** `if / else if / else`, `while`, `for` (both 3-part Go-style `for i := 0; i < N; i++` and single condition `for cond`)
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), Logical (`&&`, `||`, `!`), Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), Postfix (`++`, `--`)
- **Comments:** Single line `//` and multi-line `/* ... */`
- **Semicolons:** Optional

---

## Project Architecture

```
rustx/
├── token/       # Token definitions, TokenType constants, and keywords lookup
├── errors/      # Structured human-readable SyntaxError and RuntimeError diagnostics
├── lexer/       # Tokenizer with line/column tracking, escape decoding, and comment skipping
├── ast/         # Extensible Abstract Syntax Tree node definitions
├── parser/      # Pratt / Recursive Descent parser with precedence handling
├── object/      # Object system, values, types, and Environment scoping
├── evaluator/   # Tree-walking evaluator and standard library built-ins
├── repl/        # Interactive multiline REPL prompt
├── samples/     # Example .rx programs
├── main.go      # CLI entry point
└── go.mod       # Go module configuration (github.com/username/rustx)
```

---

## Installation and Building

### Prerequisites
- Go 1.20+ installed

### Build Binary
To build the `rustx` CLI binary:

```bash
go build -o rustx main.go
```

To run unit tests across all packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Execute Scripts Directly
Run any `.rx` script using the `run` command or shorthand syntax:

```bash
./rustx run samples/hello.rx
# or shorthand:
./rustx samples/hello.rx
```

### 2. Interactive REPL
Start the interactive shell:

```bash
./rustx repl
# or simply:
./rustx
```

Example REPL session:

```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum =", x + y)
Sum = 30
```

### 3. Build Check (Syntax Validation)
Check a `.rx` file for lexical and syntactic validity without executing:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Information

```bash
./rustx version
```

---

## Language Syntax & Features

### Variables
Variables are declared using `var` or shorthand `:=`:

```rustx
var age int = 16
name := "Muhammad"
var score = 95.5
```

### Functions & Entry Point
Functions are defined using `func`. If a zero-argument `main()` function is present in a script file, RustX automatically executes `main()` as the entry point:

```rustx
func add(a int, b int) int {
    return a + b
}

func main() {
    println("Result:", add(10, 20))
}
```

### Control Flow

#### If / Else
```rustx
if age >= 18 {
    println("Adult")
} else if age >= 13 {
    println("Teenager")
} else {
    println("Child")
}
```

#### For Loops
Supports Go-style 3-part header as well as single-condition header:

```rustx
for i := 0; i < 5; i++ {
    println("i =", i)
}

count := 3
for count > 0 {
    println("count =", count)
    count--
}
```

#### While Loops
```rustx
i := 0
while i < 3 {
    println("while loop i =", i)
    i++
}
```

### Standard Library Built-ins

RustX Starter includes built-in functions:

- **Printing:** `println(...)`, `print(...)`
- **Strings:** `len(s)`, `contains(s, sub)`, `split(s, sep)`, `join(arr, sep)`, `to_upper(s)`, `to_lower(s)`
- **Math:** `abs(n)`, `pow(x, y)`, `sin(rad)`, `cos(rad)`, `tan(rad)`
- **Time:** `sleep(ms)`, `now()`
- **Filesystem:** `read_file(path)`, `write_file(path, content)`

---

## Future Extensibility

The architecture of RustX Starter is built to accommodate future language evolution:
- **Ownership & Borrow Checker:** AST nodes and Environment scoping can be augmented with lifetime annotations and a static analysis phase prior to evaluation.
- **Bytecode VM / Native Compiler:** AST structures map directly to intermediate representations suitable for bytecode generation or LLVM lowering.
- **Concurrency & Async:** Green threads/goroutines and `async/await` syntax can be added by introducing runtime scheduler primitives to the `object` package.
