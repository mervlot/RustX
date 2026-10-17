# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

A minimal, clean, scalable interpreted/compiled hybrid programming language inspired by **Rust** and **Go**. RustX Starter blends the expressive syntax of Rust with the simplicity and speed of Go.

> **Note on Architecture:** RustX Starter currently operates as a lightweight scripting hybrid language. Memory safety systems (borrow checking, ownership model, lifetimes) are omitted in this foundation phase, but the lexer, parser, AST, and evaluator are architected so that ownership and static analysis passes can be introduced in future iterations without rewriting the compiler core.

---

## Language Highlights & Features

- **Blended Syntax:** `func` for functions, `var` for mutable variable declarations, `:=` for shorthand definitions, and optional semicolons.
- **Strong & Simple Typing:** Support for `int`, `float`, `bool`, `string`, and `array`.
- **Flexible Control Flow:** `if`, `else`, `while`, and 3-part or single-condition `for` loops with optional condition parentheses.
- **Functions & Recursion:** First-class functions with parameter and return type annotations.
- **Standard Library:** Built-in functions for printing, string manipulation (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math operations (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem basics (`read_file`, `write_file`).
- **Interactive REPL:** Multiline code entry with automatic brace/parentheses balancing.
- **Human-Readable Diagnostics:** Detailed syntax and runtime error reporting with line and column information.

---

## Directory & Architecture Overview

```text
rustx/
├── ast/            # Abstract Syntax Tree node definitions
├── errors/         # Syntax and Runtime error structures with location info
├── evaluator/      # AST evaluator and built-in function implementations
├── lexer/          # Tokenizer scanning source into token streams
├── object/         # Runtime object types and lexical Environment
├── parser/         # Pratt parser and recursive descent engine
├── repl/           # Interactive REPL with multiline balance tracking
├── samples/        # Example .rx programs
├── token/          # Token types and definitions
├── main.go         # CLI entrypoint
├── go.mod          # Go module specification
└── README.md       # Documentation
```

---

## Getting Started

### Installation & Build

Requires **Go 1.22+**.

To compile the `rustx` executable binary:

```bash
go build -o rustx main.go
```

To run all unit and integration tests across the codebase:

```bash
go test ./...
```

---

## Usage & CLI Commands

### 1. Interactive REPL

Start the REPL prompt:

```bash
./rustx repl
# or simply:
./rustx
```

Example session:
```text
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> exit
```

### 2. Executing Scripts

Run a `.rx` script directly:

```bash
./rustx run samples/hello.rx
# Or shorthand:
./rustx samples/hello.rx
```

### 3. Syntax Validation / Build

Check a `.rx` script for syntax errors without evaluating:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Info

```bash
./rustx version
```

---

## Code Examples

### 1. Hello World (`samples/hello.rx`)

```rustx
func main() {
    println("Hello, RustX Starter!")

    name := "Muhammad"
    var age int = 16
    println("Name:", name)
    println("Age:", age)
}
```

### 2. Recursion (`samples/factorial.rx`)

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    result := factorial(5)
    println("Factorial of 5 is:", result)
}
```

### 3. Control Flow (`samples/control_flow.rx`)

```rustx
func main() {
    var count = 0
    while count < 3 {
        println("Count:", count)
        count++
    }

    for i := 0; i < 5; i++ {
        if i % 2 == 0 {
            println(i, "is even")
        } else {
            println(i, "is odd")
        }
    }
}
```

---

## Future Roadmap

The modular architecture of RustX Starter is specifically designed for seamless extension in future versions:

- [ ] Ownership and Borrow Checker pass
- [ ] Explicit Lifetimes & References (`&`, `&mut`)
- [ ] Structs, Enums, and Pattern Matching (`match`)
- [ ] Traits / Interfaces
- [ ] Module and Import system (`import`)
- [ ] Async / Await and Concurrency primitives
- [ ] Bytecode VM compiler and Native JIT backend
