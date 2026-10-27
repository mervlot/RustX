# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a clean, minimal, and scalable compiled/interpreted programming language designed as a blend of **Rust** and **Go**. It combines the elegant simplicity, readability, and speed of Go with the expressive syntax and systems foundation of Rust.

> **Note on Memory Safety Architecture:**
> RustX Starter deliberately omits a borrow checker, ownership model, and lifetimes in this initial foundation phase. However, the lexer, Pratt parser, AST structures, and environment hierarchy are designed so that static analysis passes (ownership checking, lifetimes, borrow checking, and async transformations) can be attached seamlessly in future iterations without refactoring the compiler core.

---

## Language Features

- **Variables & Assignment:** Mutable variables with `var`, shorthand assignment with `:=`, and reassignments with `=`. Optional type annotations supported.
- **Data Types:** `int` (int64), `float` (float64), `bool`, `string`, and `array` (`[elem1, elem2]`).
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
  - Post-Increment / Post-Decrement: `++`, `--`
- **Control Flow:**
  - `if`, `else if`, `else`
  - `while` loops
  - `for` loops (3-part `for i := 0; i < 10; i++ { ... }` or single-condition `for condition { ... }`)
- **Functions:** Named and anonymous functions, optional parameter and return type declarations, lexical scoping, and full recursion support.
- **Comments:** Single-line `//` and multiline `/* ... */` comments.
- **Human-Readable Diagnostics:** Position-aware line and column error diagnostics (`SyntaxError` and `RuntimeError`).
- **Standard Library:** Built-in functions for strings (`str_contains`, `str_upper`, `str_lower`, `str_split`, `str_join`), math (`math_abs`, `math_sqrt`, `math_pow`, `math_floor`, `math_ceil`), time (`time_now`, `time_sleep`), and filesystem (`read_file`, `write_file`, `file_exists`).

---

## Syntax Overview

### Hello World
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Variables & Shorthand Assignment
```rustx
var age int = 16
name := "Muhammad"
println(name, "is", age)
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
    println("Factorial of 5 =", factorial(5))
}
```

### Control Flow & Loops
```rustx
for i := 0; i < 5; i++ {
    if i % 2 == 0 {
        println("Even:", i)
    } else {
        println("Odd:", i)
    }
}
```

### Arrays & Standard Library
```rustx
func main() {
    msg := "rustx starter"
    println("Upper:", str_upper(msg))

    arr := [10, 20, 30]
    arr[0] = 99
    println("Length:", len(arr), "First element:", arr[0])
}
```

---

## Project Architecture

```
rustx/
├── ast/         # AST nodes implementing Node, Statement, and Expression interfaces
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostic formatting
├── evaluator/   # Runtime evaluator, Environment scoping, and built-in function registry
├── lexer/       # Lexical analyzer with position tracking and comment skipping
├── object/      # Runtime value types (Integer, Float, Boolean, String, Array, Function, etc.)
├── parser/      # Pratt (top-down operator precedence) recursive descent parser
├── repl/        # Interactive REPL prompt with multiline input buffering
├── samples/     # Sample RustX scripts (.rx) demonstrating features
├── token/       # Token definition and keyword map
├── main.go      # CLI entry point (run, build, repl, version)
└── README.md
```

---

## CLI & REPL Usage

### Building RustX Binary
```bash
go build -o rustx main.go
```

### Commands

#### Run a Script
```bash
./rustx run samples/factorial.rx
```
or shorthand:
```bash
./rustx samples/factorial.rx
```

#### Syntax Check (Build)
Validates syntax and lexing without executing:
```bash
./rustx build samples/hello.rx
```

#### REPL (Interactive Prompt)
Launch interactive REPL mode:
```bash
./rustx repl
# or simply
./rustx
```

REPL Session Example:
```
RustX Starter v0.1
Type 'exit' or press Ctrl+C to exit.
>>> x := 5
>>> x = x * 2
>>> println(x)
10
>>> func add(a, b) {
...     return a + b
... }
>>> add(3, 4)
7
```

#### Version Info
```bash
./rustx version
```

---

## Running Tests

Run unit tests across the entire codebase:
```bash
go test -v ./...
```

---

## Future Compatibility & Extensibility

RustX Starter was architected to serve as the foundation for a full-fledged systems language. Future additions can include:
- **Ownership & Borrow Checker:** Linear/affine type checking pass over AST before evaluation.
- **Structs & Enums:** Algebraic data types and pattern matching.
- **Traits / Interfaces:** Static and dynamic polymorphism.
- **Async/Await & Concurrency:** Goroutine-like lightweight threads or futures.
- **Native / Bytecode Compiler:** LLVM or custom VM backend target reusing the AST.
