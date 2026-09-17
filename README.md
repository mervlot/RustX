# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by Rust and Go, focusing on simplicity, execution speed, readability, and future extensibility.

It combines familiar Go-style control flow and shorthand declaration syntax with Rust-like structure and keywords, operating initially as a lightweight scripting + systems hybrid language.

> **Note on Architecture:** Advanced memory safety systems (ownership, borrow checking, lifetimes) are intentionally deferred in this foundation phase. However, the lexer, parser, AST nodes, and runtime environment are architected so that type-checking, ownership tracking, bytecode compilation, or LLVM native compilation can be layered on top without structural refactoring.

---

## Language Features Overview

- **Variables & Declarations**:
  - `var` keyword with optional type annotations (`var count int = 10` or `var x = 5`).
  - `:=` shorthand assignment (`name := "RustX"`).
- **Data Types**:
  - Primitives: `int`, `float`, `bool`, `string`, `null`.
  - Composite: `array` (`[1, 2, 3]`).
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`, `%`.
  - Increment / Decrement: `++`, `--` (postfix).
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`.
  - Logical: `&&`, `||`, `!`.
- **Functions**:
  - Declared via `func name(param type) return_type { ... }`.
  - First-class anonymous functions and recursion supported.
  - Zero-argument `main()` entry-point support.
- **Control Flow**:
  - `if` / `else if` / `else` conditionals (optional parentheses).
  - `while` loops.
  - `for` loops (both 3-part `for i := 0; i < 10; i++` and single-condition `for i < 10`).
- **Standard Library Built-ins**:
  - I/O: `println(...)`, `print(...)`.
  - String operations: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`.
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`.
  - Filesystem: `read_file(path)`, `write_file(path, content)`.
- **Comments**:
  - Single line (`// comment`) and multi-line (`/* comment */`).
- **Semicolons**:
  - Optional at end of statements.

---

## Directory & Package Structure

```
.
├── ast/         # Abstract Syntax Tree interface and node definitions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime AST evaluator & standard library built-ins
├── lexer/       # Tokenizer with line/column tracking & comment handling
├── object/      # Object system, data types, & lexical environment scoping
├── parser/      # Pratt expression parser & recursive descent statement parser
├── repl/        # Interactive REPL prompt with multiline input support
├── samples/     # Example .rx script files
├── token/       # Token definitions and keyword lookup
├── main.go      # CLI entry point
└── go.mod       # Go module declaration
```

---

## Building and Running

### Prerequisites
- **Go 1.20+** installed.

### Compilation

To compile the `rustx` executable binary:

```bash
go build -o rustx main.go
```

To run unit tests across all packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Run a Script
Execute a `.rx` script directly or via the `run` command:

```bash
./rustx run samples/01_hello.rx
# OR shorthand syntax:
./rustx samples/01_hello.rx
```

### 2. Interactive REPL
Launch the interactive REPL prompt:

```bash
./rustx repl
# OR simply:
./rustx
```

Example REPL session:
```
RustX Starter v0.1
Type 'exit' or 'quit' to exit.
>>> x := 5
>>> y := 10
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(x, y)
15
```

### 3. Syntax Check / Validate Script
Validate a script's syntax without executing it:

```bash
./rustx build samples/02_factorial.rx
```

### 4. Version Info
Display the current version of RustX Starter:

```bash
./rustx version
```

---

## Example Programs

### Example 1: Hello World
`samples/01_hello.rx`
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Example 2: Recursive Factorial
`samples/02_factorial.rx`
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

### Example 3: Control Flow & Loops
`samples/03_variables_loops.rx`
```rustx
var count = 0
name := "RustX"

if name == "RustX" {
    println("Language initialized:", name)
}

println("3-part for loop:")
for i := 0; i < 5; i++ {
    println("i =", i)
}

println("While loop count down:")
x := 3
while x > 0 {
    println("x =", x)
    x--
}
```

---

## Error Handling Diagnostics

Errors in RustX Starter provide human-readable line and column positions without panicking or crashing the runner:

```
SyntaxError: Expected next token to be RPAREN, got INT instead at line 5 column 12
RuntimeError: division by zero at line 12 column 8
```

---

## Future Compatibility & Extensibility

RustX Starter is designed as a foundational stepping stone. Its clean decoupling between AST nodes, environment scoping, and evaluation enables future releases to support:
- **Borrow Checker & Ownership Systems**: Via static analysis passes on AST nodes prior to evaluation/code generation.
- **Bytecode Virtual Machine or Native Compiler**: By replacing the direct AST interpreter with an AST-to-bytecode compiler and stack VM or LLVM IR emission.
- **Structs, Enums, Traits & Modules**: By expanding token types and AST nodes without altering existing statement parsing logic.
