# RustX Starter

RustX Starter is a minimal, clean, and highly extensible programming language inspired by the elegant syntax and structural patterns of **Rust** and **Go**. Built in **Go**, the language initially behaves as a lightweight scripting and systems hybrid.

## Features

- **Variables**: `var` declarations, mutable by default, and shorthand `:=` assignments.
- **Control Flow**: Parentheses-free `if-else` blocks, `while` loops, and fully specified 3-part `for` loops (`for init; cond; post {}`).
- **Functions**: Parameterized definitions, return types, recursion, and named identifiers.
- **Primitives**: Support for Integers, Floats, Booleans, and Strings (including indexing and zero-based offsets).
- **Operators**: Includes unary prefix (`-`, `!`), postfix (`++`, `--`), standard infix math, comparisons, and short-circuiting logical (`&&`, `||`) operators.
- **Standard Library Basics**:
  - `print()`, `println()`
  - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - Time: `sleep()`, `now()`
  - Filesystem: `read_file()`, `write_file()`
- **REPL**: Interactive shell with multi-line brace and parenthesis balance checking.
- **Comments**: Supports single-line `//` and block `/* */` comments.

---

## Architectural Layout

```text
/lexer       - Lexical analysis & tokenization with full line & column position tracking
/parser      - Pratt (operator precedence) and Recursive Descent Parser
/ast         - Strong AST nodes implementing Position interfaces for diagnostics
/object      - Custom object representation and variable scoping environment
/evaluator   - Tree-walking interpreter and standard library registration
/repl        - Multiline-capable REPL shell
/token       - Lexical token definitions
/errors      - High-fidelity, human-readable SyntaxError and RuntimeError definitions
main.go      - CLI command handling and runner
```

---

## Building and Installation

Make sure you have Go installed (v1.24+ recommended).

### Compile Binary

```bash
go build -o rustx main.go
```

### CLI Usage

```bash
# Run a file directly (shorthand)
./rustx samples/hello.rx

# Or via 'run' command
./rustx run samples/factorial.rx

# Start the interactive REPL
./rustx repl
# Or simply
./rustx

# Syntactically validate a file without executing (exits with 1 on syntax error)
./rustx build samples/hello.rx

# View language version
./rustx version
```

---

## Example Program

```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is", factorial(5))
```
