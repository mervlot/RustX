# RustX Starter

**RustX Starter** is a minimal, clean, scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

It combines a Go-like lightweight syntax and interpreter architecture with Rust-inspired syntax elements and design foundations.

> **Note on Architecture & Memory Model:**
> RustX Starter is intentionally designed without an immediate borrow checker, ownership model, or lifetimes in this starter phase. However, the lexer, Pratt parser, AST structures, and runtime environment are architected so that type checking, ownership, structs, traits, modules, and native/bytecode compilation can be added in future iterations without rewriting the parser or AST core.

---

## Language Features

1. **Variables & Declarations**
   - `var` keyword with optional type annotations (`var x int = 10`)
   - Shorthand assignment `:=` operator (`name := "RustX"`)
   - Post-increment `++` and post-decrement `--` operators

2. **Data Types**
   - `int` (64-bit integer)
   - `float` (64-bit float)
   - `bool` (`true`, `false`)
   - `string` (with backslash escape sequence decoding like `\n`, `\t`, `\"`)
   - `Array` (`[1, 2, 3]`)

3. **Operators**
   - Arithmetic: `+`, `-`, `*`, `/`, `%`
   - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
   - Logical: `&&`, `||` with short-circuit evaluation

4. **Functions**
   - Top-level and first-class functions
   - Function signatures with parameter and return type annotations: `func add(a int, b int) int`
   - Full recursion support
   - Automatic execution of `main()` entry point if defined in script files

5. **Control Flow**
   - `if` / `else if` / `else` (parentheses optional)
   - `while` loops
   - 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`) and single condition `for` loops

6. **Comments**
   - Single-line comments (`//`)
   - Multiline comments (`/* ... */`)

7. **Standard Library Built-ins**
   - Printing: `println(...)`, `print(...)`
   - Strings: `len`, `to_upper`, `to_lower`, `contains`, `split`, `join`
   - Math: `abs`, `pow`, `sin`, `cos`, `tan`
   - Time: `now`, `sleep`
   - Filesystem: `read_file`, `write_file`

---

## Project Structure

```
.
├── ast/          # AST node interfaces and definitions
├── errors/       # Standardized diagnostic error types (SyntaxError, RuntimeError)
├── evaluator/    # AST Interpreter, evaluation logic, and standard library builtins
├── lexer/        # Lexical analyzer / Tokenizer
├── object/       # Runtime object types and Environment scope handling
├── parser/       # Pratt Recursive Descent Parser
├── repl/         # Interactive REPL prompt with multiline support
├── samples/      # Example .rx programs
├── token/        # Token definitions and keyword lookup table
├── go.mod        # Go module specification
├── main.go       # CLI entry point
└── README.md     # Documentation
```

---

## Building and Installation

### Prerequisites
- Go 1.20 or later

### Build Executable

```bash
go build -o rustx main.go
```

---

## Usage

### 1. Execute a Script

```bash
./rustx run samples/hello.rx
```

Or run shorthand directly with a `.rx` file path:

```bash
./rustx samples/factorial.rx
```

### 2. Launch Interactive REPL

```bash
./rustx repl
# or simply
./rustx
```

REPL Session Example:
```
RustX Starter REPL v0.1.0
Type 'exit' or 'quit' to exit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func factorial(n int) int {
...     if n <= 1 { return 1 }
...     return n * factorial(n - 1)
... }
>>> println(factorial(5))
120
```

### 3. Check Syntax (`build`)

Validates the syntax of a `.rx` file without executing it:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Info

```bash
./rustx version
```

---

## Running Unit Tests

Run all package test suites:

```bash
go test ./...
```

---

## Future Roadmap

RustX Starter is designed to serve as the foundation for:
- Ownership & Borrow Checker
- Structs, Enums, and Pattern Matching
- Interfaces / Traits
- Module & Import system
- Concurrency (async/await, channels)
- Bytecode VM / LLVM native code generator
