# RustX Starter

<p align="center">
  <img src="./assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a lightweight, clean, and scalable scripting + systems hybrid programming language blending the aesthetics and elegance of **Rust** and **Go**.

---

## 🚀 Features

- **Blended Syntax**: Combines Rust-like semantics and type annotations with Go-like clean control headers and functions (`func`, `var`, `:=`).
- **Flexible Typing**: Supports `int`, `float`, `bool`, and `string`.
- **Control Flow**: Unparenthesized `if`/`else`, `while`, and 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library**: Built-in functionality for strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`sleep`, `now`), filesystem (`read_file`, `write_file`), and formatted console output (`print`, `println`).
- **REPL**: Interactive prompt with multiline block detection and state persistence.
- **Human-Readable Diagnostics**: Precise line and column error reporting for syntax and runtime issues.
- **Extensible Architecture**: Modular AST, Parser, and Object Runtime built to seamlessly accommodate future features like an ownership model, borrow checker, structs, enums, generics, and native compilation.

---

## 📁 Architecture Overview

```text
rustx/
├── token/       # Token types, keywords, and positions
├── errors/      # Human-readable SyntaxError and RuntimeError definitions
├── lexer/       # Lexical scanner supporting comments and string escape decoding
├── ast/         # Extensible Abstract Syntax Tree node definitions
├── parser/      # Pratt recursive descent parser
├── object/      # Runtime objects and environment scoping
├── evaluator/   # AST interpreter and standard library built-ins
├── repl/        # Interactive multiline REPL interface
├── samples/     # Example .rx programs
└── main.go      # CLI entry point
```

---

## 🛠️ Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.20+

### Building the Compiler/Interpreter

From the root directory:

```bash
go build -o rustx main.go
```

### Running Scripts

Execute a RustX `.rx` script directly or via the `run` command:

```bash
./rustx run samples/hello.rx
# OR shorthand
./rustx samples/factorial.rx
```

### Starting the REPL

```bash
./rustx repl
```

Example session:
```text
RustX Starter REPL v0.1.0
Type 'exit' or press Ctrl+C to leave.

>>> x := 10
>>> func double(n int) int { return n * 2 }
>>> double(x)
20
```

### Validating Code (Build Check)

```bash
./rustx build samples/control_flow.rx
```

---

## 🔬 Sample Programs

### 1. Factorial Recursion (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### 2. Control Flow & Loops (`samples/control_flow.rx`)
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized successfully")
}

var sum = 0
for i := 1; i <= 5; i++ {
    sum = sum + i
}
println("Sum from 1 to 5 is:", sum)
```

---

## 🔮 Future Extensibility Design

RustX Starter intentionally defers borrow checking, lifetimes, and advanced memory safety to maintain a light foundation. However, every AST node retains token metadata (`GetToken() token.Token`), allowing future compiler passes (such as static analysis, ownership verification, or LLVM/C code generation) to inspect types, lifetime annotations, and memory operations without breaking compatibility.
