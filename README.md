# RustX Starter

RustX Starter is a minimal, clean, and highly scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility (e.g. adding borrow checkers or compiler targets). It features a scripting + systems hybrid syntax with a Go-inspired REPL and interpreter, implemented in Go.

---

## Features

- **Variables & Scoping:** Standard assignments (`x = y`) and Go-style shorthand declarations (`x := y`). Mutable by default and supports block-scoped scoping with lexical environments.
- **Strong & Extensible Types:** Supports primitives: integers (`int`), floats (`float`), booleans (`bool`), and double-quoted strings (`string`), as well as arrays (`[1, 2, 3]`).
- **Operators:** Standard arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `<=`, `>`, `>=`), prefix (`!`, `-`), postfix (`++`, `--`), and logical short-circuiting (`&&`, `||`).
- **Control Flow:** Go/Rust-style brace blocks, unparenthesized conditions for `if`/`else` / `else if`, `while` loops, and three-clause `for` loops (e.g. `for i := 0; i < 10; i++ { ... }`).
- **Functions:** Support for named/anonymous function signatures, optional parameter/return type annotations (e.g. `func add(a int, b int) int`), and full recursion.
- **Built-in Standard Library:**
  - Printing: `print()`, `println()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Time: `sleep()`, `now()`
  - Filesystem: `read_file()`, `write_file()`
- **Multiline Interactive REPL:** Multi-line awareness using parentheses, brackets, and braces balancing.
- **Robust Error Diagnostics:** Human-readable `SyntaxError` and `RuntimeError` tracking line and column position info. Never panics on normal user code errors.

---

## Directory Structure

```text
├── ast/         # Abstract Syntax Tree nodes
├── token/       # Token types and metadata definition
├── lexer/       # Lexical analyzer with line/col tracking and comment skipping
├── parser/      # Pratt operator precedence & recursive descent parser
├── errors/      # Position-aware Syntax and Runtime error formats
├── object/      # Runtime Object model, Environment, and scoped bindings
├── evaluator/   # Tree-walking interpreter and built-ins registration
├── repl/        # Balanced multi-line terminal prompt (REPL)
├── samples/     # RustX Starter runnable scripts
└── main.go      # CLI entry point
```

---

## Build and Installation

To build the executable from the root of the repository:

```bash
go build -o rustx main.go
```

To run all unit tests:

```bash
go test ./...
```

---

## CLI Usage

### 1. Interactive REPL
To start the interactive prompt, run:

```bash
./rustx
```
or
```bash
./rustx repl
```

**Example REPL session:**
```text
RustX Starter v0.1
>>> x := 5
>>> println(x)
5
>>> func factorial(n) { if n <= 1 { 1 } else { n * factorial(n - 1) } }
>>> println(factorial(5))
120
```

### 2. Run Scripts
To run an `.rx` file directly:

```bash
./rustx run samples/factorial.rx
```
or use the shorthand execution:
```bash
./rustx samples/factorial.rx
```

### 3. Parse and Validate Scripts (Build check)
To validate a file syntactically without executing it (exit code `1` if invalid):

```bash
./rustx build samples/factorial.rx
```

### 4. Print Version
```bash
./rustx version
```

---

## Future Extensibility

The parser and AST have been modularly designed to make future enhancements seamless:
1. **Ownership / Lifetimes:** Can be statically verified in a compiler pass after AST parsing without changing the execution model.
2. **Bytecode VM:** The current tree-walking interpreter can be easily compiled or compiled down to a bytecode instruction set for execution on a VM.
3. **Generics & Traits:** Easily integrated into function signature parameters and trait interfaces on top of standard parsing structures.
