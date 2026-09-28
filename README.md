# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, scalable programming language inspired by **Rust** and **Go**. It is focused on simplicity, speed, readability, and future extensibility.

> **Note**: This foundation version acts as a lightweight scripting and systems hybrid language. Ownership, borrow checking, and lifetime checks are intentionally omitted from this phase but the architecture (AST, symbol tables, token representation) is designed to support them in future iterations.

---

## Features

- **Blended Syntax**: `func` for functions, `var` and `:=` for variable declarations, curly brace block scoping, and optional semicolons.
- **Dynamic Primitive Types**: Integers, Floats, Booleans, Strings, Arrays, and Functions.
- **Control Flow**: `if`, `else`, `while`, and Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library Built-ins**:
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `now`, `sleep`
  - Filesystem: `read_file`, `write_file`
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Interactive REPL**: Supports multiline input, expression evaluation, and persistent execution environments across commands.
- **Clean Architecture**: Modular package layout (`token`, `lexer`, `ast`, `parser`, `object`, `evaluator`, `errors`, `repl`).

---

## Directory Layout

```text
.
├── ast/        # Abstract Syntax Tree node definitions
├── errors/     # Custom error structures (SyntaxError & RuntimeError)
├── evaluator/  # Interpreter, environment, and standard library built-ins
├── lexer/      # Lexical analyzer / tokenizer
├── object/     # Runtime object system and scoping environments
├── parser/     # Pratt recursive-descent parser
├── repl/       # Interactive REPL prompt implementation
├── samples/    # Example RustX (.rx) scripts
├── token/      # Token type definitions and keyword lookup
├── main.go     # CLI entry point
└── README.md   # Project documentation
```

---

## Installation & Build Instructions

### Prerequisites

- [Go](https://golang.org/) 1.20 or newer.

### Build Binary

Compile the `rustx` executable directly from the repository root:

```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

### 1. Run a Script
Run `.rx` files directly using the `run` command or shorthand:

```bash
./rustx run samples/hello.rx
# OR shorthand:
./rustx samples/hello.rx
```

### 2. Interactive REPL
Start the interactive REPL prompt:

```bash
./rustx repl
# OR simply:
./rustx
```

Example session:
```text
RustX Starter v0.1 REPL
Type 'exit' or Ctrl+D to quit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> exit
```

### 3. Build/Validate Script Syntax
Check for lexical and syntactic errors without executing:

```bash
./rustx build samples/factorial.rx
```

### 4. Check Version
Display version information:

```bash
./rustx version
```

---

## Language Syntax Examples

### Variables & Declarations
```rustx
var age int = 16
name := "Muhammad"
var height = 5.9
```

### Functions & Recursion
```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5)) // Outputs: 120
```

### Control Flow & Loops
```rustx
if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

for i := 0; i < 5; i++ {
    println("Count:", i)
}
```

---

## Testing

Run all unit test suites across all packages:

```bash
go test ./...
```
