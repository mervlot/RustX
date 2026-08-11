# RustX Starter

A clean, modular, and extensible compiled/interpreted programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

Initially implemented as a hybrid lightweight scripting + systems language, its architecture is designed from the ground up to support future addition of a borrow checker, ownership/lifetimes systems, modules, generics, and native compilation or bytecode VM without requiring a complete rewrite.

---

## Folder Structure

```text
/
├── main.go               # Command line interface and entry point
├── go.mod                # Go module definition
├── .gitignore            # Git ignore file
├── ast/
│   └── ast.go            # Position-aware, extensible AST Node definitions
├── errors/
│   └── errors.go         # Standard SyntaxError and RuntimeError definitions
├── lexer/
│   ├── lexer.go          # Custom Lexer with multi-line comment & escape support
│   └── lexer_test.go     # Extensive unit tests for Lexer
├── parser/
│   ├── parser.go         # Pratt + Recursive Descent parser
│   └── parser_test.go    # Extensive unit tests for Parser
├── object/
│   ├── object.go         # Object-based runtime & environment system
│   └── object_test.go    # Environment & recursive scope update tests
├── evaluator/
│   ├── evaluator.go      # AST tree-walking interpreter & short-circuit evaluation
│   ├── builtins.go       # Standard library built-ins (strings, math, filesystem, etc.)
│   └── evaluator_test.go # Extensive runtime evaluation tests
├── repl/
│   ├── repl.go           # Multiline, state-preserving REPL interactive loop
│   └── repl_test.go      # REPL multiline balance tests
└── samples/
    ├── hello.rx          # Sample: Basic hello world
    ├── factorial.rx      # Sample: Recursive factorial named function
    └── loop_math.rx      # Sample: 3-part for loop, variables, and math built-ins
```

---

## Core Features

- **Variables**: `var` keyword and shorthand `:=` assignment. Scoped strictly with block level local definitions and recursive outer lookup for reassignment.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Operators**: standard arithmetic (`+`, `-`, `*`, `/`), logical with short-circuit support (`&&`, `||`), unary (`!`, `-`), post-increment/decrement (`++`, `--`), and relational (`==`, `!=`, `<`, `>`, `<=`, `>=`).
- **Control Flow**: Block-scoped conditional blocks (`if`/`else`), loop structures (`while`), and Go-style loops (`for` with support for empty/infinite loop, conditional headers, and three-part headers `for init; condition; post {}`).
- **Semicolons**: Optional at statement and expression level.
- **Functions**: Parameter parsing, optional type annotations (both parameter and return type), closures, recursion, and named function definitions.
- **Standard Library Built-ins**:
  - Printing: `print()`, `println()` (space-separated inputs with single trailing newline).
  - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`.
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`.
  - Time: `sleep()`, `now()`.
  - Filesystem: `read_file()`, `write_file()`.

---

## How to Build

From the root directory, compile the binary using Go:

```bash
go build -o rustx main.go
```

To run all unit tests (lexer, parser, object, evaluator, and repl packages):

```bash
go test ./...
```

---

## How to Run

### 1. Interactive REPL
Start the interactive REPL prompt by calling `rustx` directly or with the `repl` subcommand:

```bash
./rustx
```

Example session:
```text
RustX Starter v0.1
Type 'exit' to quit.

>>> var x = 5;
>>> func double(n) {
...     return n * 2;
... }
>>> double(x)
10
```

### 2. Execute a Script
To execute a `.rx` file, run:

```bash
./rustx run samples/hello.rx
```

### 3. Shorthand Execution
Pass a valid path to a `.rx` file as the first argument directly:

```bash
./rustx samples/factorial.rx
```

### 4. Syntactic & Lexical Build Validation
Validate a file's syntax without executing it:

```bash
./rustx build samples/hello.rx
```

If the file contains syntax errors, it prints them to standard error and exits with code `1`. Otherwise, it returns `0`.
