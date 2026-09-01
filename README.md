# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a minimal, clean, and scalable compiled/interpreted hybrid programming language inspired by **Rust** and **Go**. It balances simplicity, speed, readability, and future extensibility.

---

## 🚀 Key Features

- **Blended Syntax**: Combines Rust and Go ergonomics (`func`, `var`, `:=`, optional semicolons, strong primitive types).
- **Control Flow**: Supports `if` / `else if` / `else`, `while` loops, and standard three-part Go-style `for` loops (`for i := 0; i < 10; i++`).
- **First-Class Functions & Recursion**: Function literals, optional type signatures, named function auto-registration, and recursive calls.
- **Built-In Standard Library**:
  - Printing: `print`, `println`
  - String manipulation: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math utilities: `abs`, `pow`, `sin`, `cos`, `tan`
  - System/Time/FS: `now`, `sleep`, `read_file`, `write_file`
- **Interactive REPL**: Interactive prompt with multiline balance tracking for braces `{}` and parentheses `()`.
- **Human-Readable Diagnostics**: Precise error messages reporting line and column positions.

---

## 🛠 Project Architecture

```
.
├── ast/         # Abstract Syntax Tree definitions
├── errors/      # Custom SyntaxError and RuntimeError diagnostics with line/column info
├── evaluator/   # Tree-walk interpreter and stdlib builtins
├── lexer/       # Lexical analyzer with line & column position tracking
├── object/      # Runtime object system and scoped Environments
├── parser/      # Pratt recursive descent parser
├── repl/        # Interactive REPL prompt with multiline input
├── samples/     # Example .rx source programs
├── token/       # Token types and definitions
└── main.go      # CLI entry point
```

> **Note on Architecture & Future Compatibility**: The architecture is intentionally designed without an aggressive borrow checker or lifetime model in this foundation phase. However, every AST node implements `GetToken() token.Token` and nodes are structured to cleanly allow future passes for ownership checking, borrow checking, bytecode compilation, generics, or async/await without requiring parser rewrites.

---

## ⚡ Quick Start

### Prerequisites
- [Go 1.22+](https://go.dev/)

### Building the Compiler/Interpreter

To build the executable binary:

```bash
go build -o rustx main.go
```

### Running Programs

Run a `.rx` script directly:

```bash
./rustx run samples/01_hello.rx
```

Or pass the `.rx` file as a direct argument shorthand:

```bash
./rustx samples/02_factorial.rx
```

### Interactive REPL

Start the REPL:

```bash
./rustx repl
# or simply
./rustx
```

Example REPL session:

```text
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to exit.
>>> x := 5
>>> y := 10
>>> println("Sum:", x + y)
Sum: 15
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(x, y)
15
```

### Validating Syntax (Build Command)

Check a `.rx` file for lexical/syntax errors without executing it:

```bash
./rustx build samples/03_features.rx
```

---

## 📜 Syntax Quick Guide

### Variables
```rustx
var age int = 16;
name := "Muhammad";
var score = 100;
```

### Functions & Recursion
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n - 1);
}

func main() {
    println("Factorial of 5:", factorial(5));
}
```

### Loops
```rustx
for i := 0; i < 5; i++ {
    println("Count:", i);
}

var x = 0;
while x < 3 {
    println("x is:", x);
    x++;
}
```

---

## 🧪 Running Tests

Execute unit tests across all packages:

```bash
go test ./...
```
