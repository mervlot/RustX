# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180">
</p>

**RustX Starter** is a minimal, clean, scalable programming language inspired by **Rust** and **Go**. It balances simplicity, speed, readability, and future extensibility.

> **Architecture Design Note**: RustX Starter is intentionally designed without a borrow checker, ownership model, or lifetimes in this initial foundation phase. The AST and evaluation model are structured so those systems can be added seamlessly in future iterations without rewriting the compiler or interpreter frontend.

---

## Key Features

- **Blended Syntax**: Integrates Rust and Go syntax styles cleanly (`func`, `var`, `:=`, optional type annotations, unparenthesized control flow).
- **Flexible Execution**: Run `.rx` scripts directly, start an interactive REPL, or build/validate files via CLI.
- **Data Types**: `int`, `float`, `bool`, `string`, `null`, and dynamic arrays `[...]`.
- **Control Flow**: `if` / `else if` / `else`, `while` loops, and Go-style 3-part `for` loops.
- **Built-in Standard Library**:
  - **I/O & Printing**: `print()`, `println()`
  - **Strings**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time**: `sleep()`, `now()`
  - **Filesystem**: `read_file()`, `write_file()`
- **Human-Readable Error Diagnostics**: Location-aware `SyntaxError` and `RuntimeError` reporting with line and column numbers.

---

## Directory Architecture

```text
rustx/
├── ast/           # AST interfaces and concrete nodes (extensible for ownership/generics)
├── errors/        # Location-aware error diagnostic structures
├── evaluator/     # Tree-walk AST evaluator & stdlib built-ins
├── lexer/         # Tokenizer supporting literals, operators, comments, line/col tracking
├── object/        # Runtime object types and lexical environment scoping
├── parser/        # Recursive descent Pratt parser with optional semicolon support
├── repl/          # Interactive REPL prompt with multiline input tracking
├── samples/       # Demonstration .rx script programs
├── token/         # TokenType definitions and keyword lookup
├── main.go        # CLI entry point
└── README.md
```

---

## Quick Start & Usage

### Build CLI

To compile the `rustx` executable binary from source:

```bash
go build -o rustx main.go
```

### Run a Script

```bash
./rustx run samples/01_hello.rx
# Or shorthand:
./rustx samples/01_hello.rx
```

### Interactive REPL

Launch the interactive REPL prompt:

```bash
./rustx repl
# Or run without arguments:
./rustx
```

```text
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func greet(name string) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### Validate / Build Syntax

Check a `.rx` file for lexical and syntax errors without executing:

```bash
./rustx build samples/02_factorial.rx
```

---

## Language Syntax Examples

### Variable Declarations

```rustx
// Shorthand assignment (:=)
name := "Muhammad"

// Explicit var statement
var age int = 16

// Mutable by default
age = 17
```

### Functions & Recursion

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

### Control Flow

```rustx
func main() {
    count := 0

    while count < 3 {
        println("Count:", count)
        count++
    }

    for i := 1; i <= 5; i++ {
        if i % 2 == 0 {
            println(i, "is even")
        }
    }
}
```

---

## Future Compatibility Roadmap

The architecture of RustX Starter is built as an extensible foundation to support future language evolutions:

1. **Ownership & Borrow Checker**: Variable tracking and scope analysis layers in evaluator / compiler.
2. **Structs, Enums & Pattern Matching**: Custom type definitions and destructuring syntax.
3. **Generics & Interfaces / Traits**: Type checking and generic method instantiation.
4. **Modules & Package Manager**: Import systems and external package resolving.
5. **Bytecode VM & Native Compilation**: IR generation / LLVM backend target.

---

## Running Tests

Run unit tests across all packages:

```bash
go test ./...
```
