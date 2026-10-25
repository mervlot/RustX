# RustX Starter

**RustX Starter** is a clean, scalable, interpreted systems/scripting hybrid programming language designed as a initial foundation inspired by **Rust** and **Go**.

It blends the elegant, modern syntax of Rust with the simplicity, speed, and readability of Go.

---

## Architecture Overview

RustX Starter is designed with modularity and future scalability in mind:

- `/token`: Defines language token types, keywords, and position tracking (line and column).
- `/errors`: Standardized, human-readable error reporting (`SyntaxError` and `RuntimeError`).
- `/lexer`: Lexical analyzer scanning source code into a stream of tokens, supporting escape characters and single/multi-line comments.
- `/ast`: Abstract Syntax Tree data structures designed to be extended for future language constructs ( ownership, borrow checker, structs, generics, etc.).
- `/parser`: Pratt recursive descent parser handling optional semicolons, expression priorities, and syntax verification.
- `/object`: Object runtime environment supporting dynamic object types and lexical scope management.
- `/evaluator`: Tree-walking interpreter evaluating expressions, control flow, functions, and standard library built-ins.
- `/repl`: Interactive REPL prompt supporting multiline input and session persistent environment.
- `/samples`: Collection of `.rx` sample programs demonstrating language features.

---

## Features

- **Variables & Assignment:** `var` declarations, `:=` shorthand assignment.
- **Data Types:** Integers, Floats, Booleans, Strings, Arrays.
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
  - Postfix: `++`, `--`
- **Control Flow:** `if` / `else if` / `else`, `while` loops, C/Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Functions:** Parameters with optional type annotations, return types, recursive call support, and automatic entry-point execution via `main()`.
- **Standard Library:**
  - I/O: `print()`, `println()`
  - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - System/Time: `now()`, `sleep()`
  - Filesystem: `read_file()`, `write_file()`

---

## Building and Installation

Requires **Go 1.20+**.

```bash
# Build the rustx binary
go build -o rustx main.go
```

---

## Usage

### Run a Script

```bash
./rustx run samples/hello.rx
# OR
./rustx samples/hello.rx
```

### Start the REPL

```bash
./rustx repl
# OR simply
./rustx
```

Example REPL interaction:
```text
RustX Starter v0.1
Type 'exit' or 'quit' to exit the REPL.
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

### Validate Syntax (Build Command)

```bash
./rustx build samples/factorial.rx
```

---

## Testing

Run all unit and integration test suites:

```bash
go test ./...
```
