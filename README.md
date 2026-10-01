# RustX Starter

RustX Starter is a lightweight, clean, and scalable compiled/interpreted programming language inspired by **Rust** and **Go**. Designed with simplicity, speed, readability, and future extensibility in mind, it provides a scripting + systems hybrid foundation.

> **Note:** The current foundation phase deliberately omits an active borrow checker, ownership model, or lifetime system to keep the language lightweight and approachable. The architecture is explicitly designed so memory safety and systems features can be integrated seamlessly in future versions without rewriting core components.

---

## Language Features

- **Blended Rust + Go Syntax**: `func` for functions, `var` for declarations, `:=` shorthand assignment, optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, arrays (`[...]`).
- **Control Flow**: `if` / `else if` / `else`, `while`, and 3-part or conditional `for` loops.
- **Functions & Recursion**: Named function literals, recursive calls, and parameter support.
- **Comments**: Single-line `//` and multi-line `/* */`.
- **Standard Library**:
  - Printing: `print`, `println`
  - Strings: `len`, `to_upper`, `to_lower`, `contains`, `split`, `join`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `now`, `sleep`
  - Filesystem: `read_file`, `write_file`

---

## Directory & Architecture Overview

```
rustx/
├── ast/         # Abstract Syntax Tree nodes for statements & expressions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Evaluator engine and standard library builtins
├── lexer/       # Lexical analyzer with position tracking and escape sequence decoding
├── object/      # Runtime value objects and lexical environment hierarchy
├── parser/      # Pratt parser implementation with error recovery
├── repl/        # Interactive REPL with multiline block tracking
├── samples/     # Example programs (.rx files)
├── token/       # Token types, position structures, and lookup maps
├── main.go      # CLI entry point (run, repl, build, version)
└── README.md
```

---

## Building and Installation

### Requirements
- Go 1.20 or newer

### Build Executable
```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

### 1. Execute Scripts
```bash
./rustx run samples/hello.rx
# Or shorthand execution:
./rustx samples/hello.rx
```

### 2. Interactive REPL
```bash
./rustx repl
# Or run without args:
./rustx
```

Example REPL session:
```
RustX Starter v0.1.0 Interactive REPL
Type statements or expressions. Press Ctrl+C or Ctrl+D to exit.

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

### 3. Syntax Validation / Build Check
```bash
./rustx build samples/demo.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## Syntax Examples

### Example 1: Hello World
```rustx
func main() {
    println("Hello World")
}
```

### Example 2: Recursive Factorial
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

### Example 3: Loop & Standard Library
```rustx
func main() {
    name := "RustX"
    if name == "RustX" {
        println("Language initialized")
    }

    for i := 0; i < 5; i++ {
        println("Count:", i)
    }

    msg := "rustx starter"
    println(to_upper(msg))
}
```

---

## Testing

Run the full test suite across all packages:
```bash
go test ./...
```
