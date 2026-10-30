# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a clean, scalable programming language inspired by **Rust** and **Go**. It blends elegant syntax with high-performance execution, combining a scripting feel with systems-level design foundations.

> **Note**: RustX Starter is designed as a clean foundation. Advanced features like ownership models, borrow checking, async/await, traits, and modules are architected to be added in future iterations without rewriting the compiler core.

---

## Features

- **Blended Rust & Go Syntax**: `func` for functions, `var` and `:=` shorthand for variable declaration, curly brace blocks, optional semicolons.
- **Dynamic Primitive Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Operators**: Full arithmetic (`+`, `-`, `*`, `/`, `%`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), and postfix (`++`, `--`).
- **Control Flow**:
  - `if` / `else if` / `else`
  - `while` loops
  - `for` loops (supports both 3-part `for i := 0; i < 10; i++` and single-condition `for i < 10`).
- **Standard Library Built-ins**:
  - **Console**: `println(...)`, `print(...)`
  - **Containers & Strings**: `len(val)`, `str_to_upper(s)`, `str_to_lower(s)`, `str_contains(s, sub)`
  - **Math**: `sqrt(x)`, `pow(x, y)`
  - **Time**: `now_ms()`, `sleep_ms(ms)`
  - **Filesystem**: `read_file(path)`, `write_file(path, content)`
- **Interactive REPL**: Multiline code support with state persistence across entries.
- **Human-Readable Diagnostics**: Detailed line and column reporting for syntax and runtime errors.

---

## Folder Architecture

```text
.
├── ast/         # Abstract Syntax Tree node definitions
├── errors/      # Human-readable SyntaxError & RuntimeError types
├── evaluator/   # Runtime interpreter and standard library built-ins
├── lexer/       # Tokenizer with line & column tracking
├── object/      # Internal runtime object system and lexical environments
├── parser/      # Pratt & Recursive Descent parser
├── repl/        # Interactive REPL prompt with multiline input handling
├── samples/     # Sample RustX (.rx) scripts
├── token/       # Token definitions and keyword lookup table
├── go.mod       # Go module declaration
├── main.go      # CLI entrypoint (run, repl, build, version)
└── README.md    # Language documentation
```

---

## Quick Start

### Build Executable
```bash
go build -o rustx main.go
```

### Run a Script
```bash
./rustx run samples/hello.rx
# or directly pass the filename:
./rustx samples/factorial.rx
```

### Start the REPL
```bash
./rustx repl
# or run without arguments:
./rustx
```

### Check Syntax (`build`)
```bash
./rustx build samples/control_flow.rx
```

### Version Info
```bash
./rustx version
```

---

## Language Examples

### 1. Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("RustX Starter")
}
```

### 2. Recursion & Math (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### 3. Control Flow (`samples/control_flow.rx`)
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

for i := 0; i < 3; i++ {
    println("For loop index:", i)
}
```

---

## Future Compatibility & Extensibility Roadmap

1. **Memory Model**: Incorporating static borrow checker and lifetime analysis atop the current AST.
2. **Modules & Imports**: Package management and multi-file namespace resolution.
3. **Concurrency**: Lightweight async/await tasks or Go-like channels.
4. **Custom Types**: `struct`, `enum`, `pattern matching`, and `traits`.
