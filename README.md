# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**, focused on simplicity, speed, readability, and future extensibility. It bridges the gap between lightweight scripting productivity and systems language aesthetics.

---

## Language Highlights & Syntax Style

RustX blends the elegance of Go and the expressiveness of Rust:

- **Functions**: `func name(arg type) returnType { ... }`
- **Variables**: `var x = 5` or `x := 5` (mutable by default)
- **Control Flow**: `if / else`, `while`, and `for` (both 3-part Go/C-style `for i := 0; i < 10; i++` and conditional loops)
- **Parentheses optional** in condition headers (`if x > 10 { ... }`)
- **Semicolons optional**
- **Types**: `int`, `float`, `bool`, `string`, `array`
- **Comments**: Single-line (`//`) and multi-line (`/* ... */`)
- **Extensible Architecture**: Modular architecture ready for future ownership/borrow checkers, native code generation, or bytecode VM execution without full compiler rewrites.

---

## Getting Started & Build Instructions

### Prerequisites
- Go 1.20 or later

### Building the CLI
Build the `rustx` executable binary from the root directory:

```bash
go build -o rustx main.go
```

To run all unit tests:
```bash
go test -v ./...
```

---

## CLI Usage

```text
RustX Starter - Lightweight Scripting & Systems Hybrid Language

Usage:
  rustx run <file.rx>     Execute a RustX source file
  rustx build <file.rx>   Check syntax and validate file without executing
  rustx repl              Start interactive REPL session
  rustx version           Display version information
  rustx <file.rx>         Shorthand to execute a file
  rustx                   Start REPL session
```

### Interactive REPL

Launch REPL:
```bash
./rustx repl
# or simply
./rustx
```

Example session:
```text
RustX Starter v0.1
Type 'exit' or Ctrl+D to quit.

>>> name := "Muhammad"
>>> println("Hello", name)
Hello Muhammad
>>> func add(a, b) { return a + b }
>>> add(10, 25)
35
>>> exit
```

Multiline input (block balancing with `{}` and `()`) is fully supported.

---

## Example Code

### 1. Hello World (`samples/hello.rx`)
```go
func main() {
    println("Hello World")
}
```

### 2. Recursion (`samples/factorial.rx`)
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5 is:")
    println(factorial(5))
}
```

### 3. Loops and Control Flow (`samples/loops.rx`)
```go
func main() {
    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    count := 3
    while count > 0 {
        println("count down:", count)
        count--
    }
}
```

### 4. Built-in Standard Library (`samples/stdlib.rx`)
```go
func main() {
    println("abs(-15) =", abs(-15))
    println("sqrt(16) =", sqrt(16))
    println("pow(2, 8) =", pow(2, 8))

    str := "  rustx starter  "
    println(to_upper(trim(str)))

    write_file("test.txt", "RustX file test")
    content := read_file("test.txt")
    println("File contents:", content)
}
```

---

## Standard Library Functions

| Category | Function | Description |
|----------|----------|-------------|
| **I/O** | `print(...)`, `println(...)` | Output arguments with space separation |
| **Inspection** | `len(item)`, `type(item)` | Length of strings/arrays, runtime type inspection |
| **Math** | `abs(x)` | Absolute value |
| | `sqrt(x)` | Square root |
| | `pow(base, exp)` | Power function |
| | `floor(x)`, `ceil(x)` | Rounding functions |
| | `min(a, b, ...)`, `max(a, b, ...)` | Minimum and maximum across arguments |
| **Strings** | `to_upper(str)`, `to_lower(str)` | Case conversion |
| | `trim(str)` | Strip leading/trailing whitespace |
| | `contains(str, substr)` | Substring search |
| | `split(str, sep)` | Split into string array |
| | `join(arr, sep)` | Join array into string |
| **Time** | `now()` | Unix epoch timestamp (seconds) |
| | `sleep(ms)` | Pause execution for given milliseconds |
| **Filesystem** | `read_file(path)` | Read file content into string |
| | `write_file(path, content)` | Write string content to disk |

---

## Architecture

The codebase is organized into modular packages:

```
├── ast/        # Abstract Syntax Tree nodes with source position tracking
├── errors/     # Standardized SyntaxError and RuntimeError diagnostic representations
├── evaluator/  # Tree-walking interpreter and built-in standard library
├── lexer/      # Lexical analyzer / tokenizer with comment & escape handling
├── object/     # Runtime object system and lexical environment scoping
├── parser/     # Recursive descent Pratt parser with operator precedence
├── repl/       # Interactive REPL with multiline balance support
├── samples/    # Sample .rx demonstration scripts
├── token/      # Token types, definitions, and keywords
└── main.go     # CLI entry point
```

### Future Extensibility
The AST design and environment model decouple syntactic representation from execution semantics:
- **Ownership & Lifetimes**: AST nodes capture full position and identifier relationships, allowing static analysis passes (e.g. borrow checking) to run between parsing and interpretation/compilation.
- **Bytecode VM or Native Compiler**: The AST can directly serve as the IR for emitting bytecode instructions or LLVM IR in future phases.
- **Type Checking**: Function signatures and variable annotations (`int`, `float`, etc.) are retained in the AST, facilitating static type checker passes without breaking existing dynamic runtime features.
